package main

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/darabuchi/prism"
	"github.com/darabuchi/prism/pkg/geoip"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/xerror"
	"github.com/lazygophers/utils/candy"
	"github.com/oschwald/maxminddb-golang"
)

// Converter 数据转换器
type Converter struct {
	writer   *geoip.Writer
	inserter *BatchInserter
}

// NewConverter 创建转换器
func NewConverter() (*Converter, error) {
	writer, err := geoip.NewWriter()
	if err != nil {
		log.Errorf("创建转换器失败: %v", err)
		return nil, xerror.WrapError(err, ErrWriterCreateFailed, "create converter failed")
	}

	log.Debugf("转换器创建成功")

	// 创建批量插入器，批大小为 5000
	inserter := NewBatchInserter(writer, 5000)

	return &Converter{
		writer:   writer,
		inserter: inserter,
	}, nil
}

// Flush 刷新批量插入缓存
func (c *Converter) Flush() error {
	return c.inserter.Flush()
}

// ConvertMaxMindDB 从 MaxMind 数据库转换
func (c *Converter) ConvertMaxMindDB(dbPath string) error {
	log.Infof("开始转换 MaxMind 数据库: %s", filepath.Base(dbPath))

	db, err := maxminddb.Open(dbPath)
	if err != nil {
		log.Errorf("打开 MaxMind 数据库失败，路径: %s, 错误: %v", dbPath, err)
		return xerror.WrapError(err, ErrDatabaseOpenFailed, "open maxmind db "+filepath.Base(dbPath)+" failed")
	}
	defer db.Close()

	// MaxMind 数据库通常使用网络迭代器
	networks := db.Networks(maxminddb.SkipAliasedNetworks)

	count := 0
	skipCount := 0
	errorCount := 0

	for networks.Next() {
		var record map[string]interface{}
		subnet, err := networks.Network(&record)
		if err != nil {
			errorCount++
			log.Warnf("读取网络记录失败 (第 %d 条): %v", count+skipCount+errorCount, err)
			continue
		}

		// 转换为 GeoIP
		geo := convertMaxMindRecordToGeoIP(record, subnet.IP.String())
		if geo == nil {
			skipCount++
			continue
		}

		// 获取 IP 范围
		firstIP := subnet.IP
		lastIP := getLastIP(subnet)

		// 使用批量插入
		err = c.inserter.Add(firstIP.String(), lastIP.String(), geo)
		if err != nil {
			errorCount++
			log.Errorf("插入 IP 范围失败 (第 %d 条)，范围: %s-%s, 错误: %v",
				count+skipCount+errorCount, firstIP.String(), lastIP.String(), err)
			return xerror.WrapError(err, ErrInsertFailed,
				"insert IP range "+firstIP.String()+"-"+lastIP.String()+" from "+filepath.Base(dbPath)+" failed")
		}

		count++
		if count%10000 == 0 {
			log.Infof("已处理 %d 条记录 (跳过 %d, 错误 %d)", count, skipCount, errorCount)
		}
	}

	if networks.Err() != nil {
		log.Errorf("遍历网络记录失败，数据库: %s, 错误: %v", filepath.Base(dbPath), networks.Err())
		return xerror.WrapError(networks.Err(), ErrConversionFailed,
			"iterate networks in "+filepath.Base(dbPath)+" failed")
	}

	// 刷新剩余批次
	if err := c.inserter.Flush(); err != nil {
		return err
	}

	log.Infof("完成转换 %s: 成功 %d 条，跳过 %d 条，错误 %d 条",
		filepath.Base(dbPath), count, skipCount, errorCount)
	return nil
}

// ConvertCSV 从 CSV 文件转换
func (c *Converter) ConvertCSV(csvPath string, converter func(record []string) (*prism.GeoIP, string, string, error)) error {
	log.Infof("开始转换 CSV 文件: %s", filepath.Base(csvPath))

	file, err := os.Open(csvPath)
	if err != nil {
		log.Errorf("打开 CSV 文件失败，路径: %s, 错误: %v", csvPath, err)
		return xerror.WrapError(err, ErrFileOperationFailed, "open csv file "+filepath.Base(csvPath)+" failed")
	}
	defer file.Close()

	var reader *csv.Reader
	if strings.HasSuffix(csvPath, ".gz") {
		gr, err := gzip.NewReader(file)
		if err != nil {
			log.Errorf("解压 GZ 文件失败，路径: %s, 错误: %v", csvPath, err)
			return xerror.WrapError(err, ErrFileOperationFailed, "create gzip reader for "+filepath.Base(csvPath)+" failed")
		}
		defer gr.Close()
		reader = csv.NewReader(gr)
	} else {
		reader = csv.NewReader(file)
	}

	count := 0
	skipCount := 0
	errorCount := 0
	lineNumber := 0

	for {
		record, err := reader.Read()
		lineNumber++

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("读取 CSV 行失败，文件: %s, 行号: %d, 错误: %v",
				filepath.Base(csvPath), lineNumber, err)
			return xerror.WrapError(err, ErrFileOperationFailed,
				fmt.Sprintf("read csv line %d from %s failed", lineNumber, filepath.Base(csvPath)))
		}

		// 跳过表头
		if count == 0 && len(record) > 0 && record[0] == "start_ip" {
			count++
			continue
		}

		geo, startIP, endIP, err := converter(record)
		if err != nil {
			errorCount++
			log.Warnf("转换 CSV 记录失败，文件: %s, 行号: %d, 错误: %v, 记录: %v",
				filepath.Base(csvPath), lineNumber, err, record)
			continue
		}

		if geo == nil {
			skipCount++
			continue
		}

		err = c.inserter.Add(startIP, endIP, geo)
		if err != nil {
			errorCount++
			log.Errorf("插入 CSV 数据失败，文件: %s, 行号: %d, 范围: %s-%s, 错误: %v",
				filepath.Base(csvPath), lineNumber, startIP, endIP, err)
			return xerror.WrapError(err, ErrInsertFailed,
				fmt.Sprintf("insert IP range %s-%s from %s line %d failed",
					startIP, endIP, filepath.Base(csvPath), lineNumber))
		}

		count++
		if count%10000 == 0 {
			log.Infof("已处理 %d 条记录 (跳过 %d, 错误 %d)", count, skipCount, errorCount)
		}
	}

	// 刷新剩余批次
	if err := c.inserter.Flush(); err != nil {
		return err
	}

	log.Infof("完成转换 %s: 成功 %d 条，跳过 %d 条，错误 %d 条",
		filepath.Base(csvPath), count, skipCount, errorCount)
	return nil
}

// WriteTo 写入数据库文件
func (c *Converter) WriteTo(w io.Writer) (int64, error) {
	return c.writer.WriteTo(w)
}

// convertMaxMindRecordToGeoIP 将 MaxMind 记录转换为 GeoIP
func convertMaxMindRecordToGeoIP(record map[string]interface{}, ip string) *prism.GeoIP {
	geo := &prism.GeoIP{
		IP: ip,
	}

	// Country
	if country, ok := record["country"].(map[string]interface{}); ok {
		if names, ok := country["names"].(map[string]interface{}); ok {
			if en, ok := names["en"].(string); ok {
				geo.Country = en
			}
		}
		if code, ok := country["iso_code"].(string); ok {
			geo.CountryCode = code
		}
	}

	// Continent
	if continent, ok := record["continent"].(map[string]interface{}); ok {
		if code, ok := continent["code"].(string); ok {
			geo.ContinentCode = code
		}
	}

	// Subdivisions (Region)
	if subdivisions, ok := record["subdivisions"].([]interface{}); ok && len(subdivisions) > 0 {
		if sub, ok := subdivisions[0].(map[string]interface{}); ok {
			if names, ok := sub["names"].(map[string]interface{}); ok {
				if en, ok := names["en"].(string); ok {
					geo.Region = en
				}
			}
		}
	}

	// City
	if city, ok := record["city"].(map[string]interface{}); ok {
		if names, ok := city["names"].(map[string]interface{}); ok {
			if en, ok := names["en"].(string); ok {
				geo.City = en
			}
		}
	}

	// Location
	if location, ok := record["location"].(map[string]interface{}); ok {
		if lat, ok := location["latitude"].(float64); ok {
			geo.Latitude = lat
		}
		if lon, ok := location["longitude"].(float64); ok {
			geo.Longitude = lon
		}
		if tz, ok := location["time_zone"].(string); ok {
			geo.Timezone = tz
		}
	}

	// Traits (ASN)
	if traits, ok := record["traits"].(map[string]interface{}); ok {
		if asn, ok := traits["autonomous_system_number"].(uint); ok {
			geo.ASN = int(asn)
		}
		if asOrg, ok := traits["autonomous_system_organization"].(string); ok {
			geo.ASName = asOrg
			geo.Org = asOrg
		}
		if isp, ok := traits["isp"].(string); ok {
			geo.ISP = isp
		}
	}

	// 填充派生字段
	geo.FillDerivedFields()

	return geo
}

// getLastIP 获取子网的最后一个 IP
func getLastIP(subnet *net.IPNet) net.IP {
	ip := subnet.IP
	mask := subnet.Mask

	lastIP := make(net.IP, len(ip))
	copy(lastIP, ip)

	for i := range lastIP {
		lastIP[i] |= ^mask[i]
	}

	return lastIP
}

// ConvertIPInfoCSV IPInfo CSV 格式转换器
// CSV 格式: start_ip,end_ip,country,country_name,continent,continent_name,asn,as_name,as_domain
// 参考: https://ipinfo.io/products/free-ip-database
func ConvertIPInfoCSV(record []string) (*prism.GeoIP, string, string, error) {
	// IPInfo CSV 至少需要 8 个字段
	if len(record) < 8 {
		return nil, "", "", fmt.Errorf("invalid record length: %d, expected at least 8", len(record))
	}

	// 解析 ASN，去除 "AS" 或 "as" 前缀
	asnStr := strings.ToUpper(strings.TrimSpace(record[6]))
	asnStr = strings.TrimPrefix(asnStr, "AS")
	asnNum := int(candy.ToUint64(asnStr))

	geo := &prism.GeoIP{
		CountryCode:   strings.TrimSpace(record[2]),
		ContinentCode: strings.TrimSpace(record[4]),
		ASN:           asnNum,
		ASName:        strings.TrimSpace(record[7]),
		Org:           strings.TrimSpace(record[7]),
	}

	// 如果有 ASN，设置 AS 字符串表示（如 "AS15169"）
	if geo.ASN > 0 {
		geo.AS = fmt.Sprintf("AS%d", geo.ASN)
	}

	// 填充派生字段（国家名称等）
	geo.FillDerivedFields()

	return geo, strings.TrimSpace(record[0]), strings.TrimSpace(record[1]), nil
}

// ConvertDBIPCountryCSV DBIP Country CSV 格式转换器
func ConvertDBIPCountryCSV(record []string) (*prism.GeoIP, string, string, error) {
	if len(record) < 3 {
		return nil, "", "", fmt.Errorf("invalid record length: %d", len(record))
	}

	geo := &prism.GeoIP{
		CountryCode: record[2],
	}

	return geo, record[0], record[1], nil
}

// ConvertDBIPASNCSV DBIP ASN CSV 格式转换器
func ConvertDBIPASNCSV(record []string) (*prism.GeoIP, string, string, error) {
	if len(record) < 4 {
		return nil, "", "", fmt.Errorf("invalid record length: %d", len(record))
	}

	asn, _ := strconv.Atoi(record[2])
	geo := &prism.GeoIP{
		ASN:    asn,
		ASName: record[3],
		Org:    record[3],
	}

	geo.FillDerivedFields()

	return geo, record[0], record[1], nil
}

// ConvertDBIPCityCSV DBIP City CSV 格式转换器
func ConvertDBIPCityCSV(record []string) (*prism.GeoIP, string, string, error) {
	if len(record) < 10 {
		return nil, "", "", fmt.Errorf("invalid record length: %d", len(record))
	}

	lat, _ := strconv.ParseFloat(record[7], 64)
	lon, _ := strconv.ParseFloat(record[8], 64)

	geo := &prism.GeoIP{
		CountryCode: record[2],
		City:        record[5],
		Latitude:    lat,
		Longitude:   lon,
		Timezone:    record[9],
	}

	return geo, record[0], record[1], nil
}

// ConvertChinaIPListTXT China IP List TXT 格式转换器
func ConvertChinaIPListTXT(cidr string) (*prism.GeoIP, string, string, error) {
	cidr = strings.TrimSpace(cidr)
	if cidr == "" || strings.HasPrefix(cidr, "#") {
		return nil, "", "", nil
	}

	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Warnf("解析 CIDR 失败: %s, 错误: %v", cidr, err)
		return nil, "", "", xerror.WrapError(err, ErrInvalidIP, "parse CIDR "+cidr+" failed")
	}

	firstIP := ip.Mask(ipNet.Mask)
	lastIP := getLastIP(ipNet)

	geo := &prism.GeoIP{
		CountryCode: "CN",
		Country:     "China",
	}

	return geo, firstIP.String(), lastIP.String(), nil
}
