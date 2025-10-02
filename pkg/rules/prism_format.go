package rules

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"
	"net/netip"
	"strings"

	"github.com/darabuchi/prism"
)

const (
	// MagicNumber Prism 二进制文件魔数
	MagicNumber = "PRISM"

	// Version 当前格式版本
	Version = uint8(2) // v2: gzip + 优化编码

	// HeaderSize 文件头大小（字节）
	HeaderSize = 16
)

// BinaryRuleType 二进制规则类型编码
type BinaryRuleType uint8

const (
	BinaryTypeDomain       BinaryRuleType = 1
	BinaryTypeDomainSuffix BinaryRuleType = 2
	BinaryTypeIPCIDR       BinaryRuleType = 3
	BinaryTypeIPCIDR6      BinaryRuleType = 4
)

// 常见域名后缀字典（按频率排序）
var domainSuffixDict = []string{
	".com",
	".net",
	".org",
	".io",
	".cn",
	".cloudfront.net",
	".akamaized.net",
	".akamaihd.net",
	".amazonaws.com",
	".cloudflare.com",
	".googleapis.com",
	".azure.com",
	".windows.net",
	".apple.com",
	".google.com",
	".facebook.com",
	".microsoft.com",
	".co.uk",
	".de",
	".fr",
	".jp",
	".kr",
	".tw",
	".hk",
	".sg",
}

// PrismBinaryWriter Prism 二进制格式写入器
type PrismBinaryWriter struct {
	w              io.Writer
	gzipEnabled    bool
	gzWriter       *gzip.Writer
	suffixDict     map[string]uint8
	reverseSuffix  []string
	varintEnabled  bool
	ipBinaryEnabled bool
}

// NewPrismBinaryWriter 创建新的 Prism 二进制写入器
func NewPrismBinaryWriter(w io.Writer, opts ...BinaryOption) *PrismBinaryWriter {
	writer := &PrismBinaryWriter{
		w:              w,
		gzipEnabled:    true,  // 默认启用 gzip
		varintEnabled:  true,  // 默认启用变长整数
		ipBinaryEnabled: true, // 默认启用 IP 二进制编码
	}

	// 应用选项
	for _, opt := range opts {
		opt(writer)
	}

	// 构建后缀字典
	writer.buildSuffixDict()

	return writer
}

// BinaryOption 二进制编码选项
type BinaryOption func(*PrismBinaryWriter)

// WithGzip 启用/禁用 gzip 压缩
func WithGzip(enabled bool) BinaryOption {
	return func(w *PrismBinaryWriter) {
		w.gzipEnabled = enabled
	}
}

// WithVarint 启用/禁用变长整数编码
func WithVarint(enabled bool) BinaryOption {
	return func(w *PrismBinaryWriter) {
		w.varintEnabled = enabled
	}
}

// WithIPBinary 启用/禁用 IP 二进制编码
func WithIPBinary(enabled bool) BinaryOption {
	return func(w *PrismBinaryWriter) {
		w.ipBinaryEnabled = enabled
	}
}

// buildSuffixDict 构建域名后缀字典
func (w *PrismBinaryWriter) buildSuffixDict() {
	w.suffixDict = make(map[string]uint8, len(domainSuffixDict))
	w.reverseSuffix = make([]string, len(domainSuffixDict))

	for i, suffix := range domainSuffixDict {
		w.suffixDict[suffix] = uint8(i)
		w.reverseSuffix[i] = suffix
	}
}

// compressDomain 压缩域名（提取常见后缀）
func (w *PrismBinaryWriter) compressDomain(domain string) (compressed string, suffixID uint8) {
	for suffix, id := range w.suffixDict {
		if strings.HasSuffix(domain, suffix) {
			return strings.TrimSuffix(domain, suffix), id
		}
	}
	return domain, 255 // 255 表示无匹配后缀
}

// decompressDomain 解压域名
func (w *PrismBinaryWriter) decompressDomain(compressed string, suffixID uint8) string {
	if suffixID == 255 {
		return compressed
	}
	if int(suffixID) < len(w.reverseSuffix) {
		return compressed + w.reverseSuffix[suffixID]
	}
	return compressed
}

// writeVarint 写入变长整数
func (w *PrismBinaryWriter) writeVarint(writer io.Writer, value uint64) error {
	if !w.varintEnabled {
		return binary.Write(writer, binary.LittleEndian, uint32(value))
	}

	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, value)
	_, err := writer.Write(buf[:n])
	return err
}

// readVarint 读取变长整数
func readVarint(reader io.Reader, varintEnabled bool) (uint64, error) {
	if !varintEnabled {
		var val uint32
		err := binary.Read(reader, binary.LittleEndian, &val)
		return uint64(val), err
	}

	return binary.ReadUvarint(byteReader{reader})
}

// byteReader 包装 io.Reader 为 io.ByteReader
type byteReader struct {
	io.Reader
}

func (br byteReader) ReadByte() (byte, error) {
	var b [1]byte
	_, err := br.Reader.Read(b[:])
	return b[0], err
}

// Write 写入规则集到二进制格式
func (w *PrismBinaryWriter) Write(rulesByAction map[string][]Rule) error {
	// 写入文件头（不压缩）
	if err := w.writeHeader(w.w); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	var writer io.Writer = w.w

	// 如果启用 gzip，创建 gzip writer（只压缩数据部分）
	if w.gzipEnabled {
		w.gzWriter = gzip.NewWriter(w.w)
		writer = w.gzWriter
		defer w.gzWriter.Close()
	}

	// 写入规则数据
	if err := w.writeRules(writer, rulesByAction); err != nil {
		return fmt.Errorf("write rules: %w", err)
	}

	return nil
}

// writeHeader 写入文件头
func (w *PrismBinaryWriter) writeHeader(writer io.Writer) error {
	// Magic Number (5 bytes)
	if _, err := writer.Write([]byte(MagicNumber)); err != nil {
		return err
	}

	// Version (1 byte)
	if err := binary.Write(writer, binary.LittleEndian, Version); err != nil {
		return err
	}

	// Flags (1 byte)
	var flags uint8
	if w.gzipEnabled {
		flags |= 0x01
	}
	if w.varintEnabled {
		flags |= 0x02
	}
	if w.ipBinaryEnabled {
		flags |= 0x04
	}
	if err := binary.Write(writer, binary.LittleEndian, flags); err != nil {
		return err
	}

	// Reserved (9 bytes)
	reserved := make([]byte, 9)
	if _, err := writer.Write(reserved); err != nil {
		return err
	}

	return nil
}

// writeRules 写入规则数据
func (w *PrismBinaryWriter) writeRules(writer io.Writer, rulesByAction map[string][]Rule) error {
	// 写入 Action 数量
	if err := binary.Write(writer, binary.LittleEndian, uint16(len(rulesByAction))); err != nil {
		return err
	}

	// 写入每个 Action 及其规则
	for action, ruleList := range rulesByAction {
		if err := w.writeAction(writer, action, ruleList); err != nil {
			return fmt.Errorf("write action %s: %w", action, err)
		}
	}

	return nil
}

// writeAction 写入单个 Action 及其规则
func (w *PrismBinaryWriter) writeAction(writer io.Writer, action string, ruleList []Rule) error {
	// 写入 Action 名称长度
	if err := binary.Write(writer, binary.LittleEndian, uint8(len(action))); err != nil {
		return err
	}

	// 写入 Action 名称
	if _, err := writer.Write([]byte(action)); err != nil {
		return err
	}

	// 过滤并转换规则
	validRules := w.filterRules(ruleList)

	// 写入规则数量（使用 varint）
	if err := w.writeVarint(writer, uint64(len(validRules))); err != nil {
		return err
	}

	// 写入每条规则
	for _, rule := range validRules {
		if err := w.writeRule(writer, rule); err != nil {
			return fmt.Errorf("write rule: %w", err)
		}
	}

	return nil
}

// filterRules 过滤并转换规则（只保留支持的类型）
func (w *PrismBinaryWriter) filterRules(ruleList []Rule) []binaryRule {
	validRules := make([]binaryRule, 0, len(ruleList))

	for _, rule := range ruleList {
		ruleType := rule.Type()
		payload := rule.Payload()

		var bt BinaryRuleType
		var ok bool

		switch ruleType {
		case TypeDomain:
			bt = BinaryTypeDomain
			ok = true
		case TypeDomainSuffix:
			bt = BinaryTypeDomainSuffix
			ok = true
		case TypeIPCIDR:
			bt = BinaryTypeIPCIDR
			ok = true
		case TypeIPCIDR6:
			bt = BinaryTypeIPCIDR6
			ok = true
		default:
			continue
		}

		if ok {
			validRules = append(validRules, binaryRule{
				ruleType: bt,
				payload:  payload,
			})
		}
	}

	return validRules
}

type binaryRule struct {
	ruleType BinaryRuleType
	payload  string
}

// writeRule 写入单条规则
func (w *PrismBinaryWriter) writeRule(writer io.Writer, rule binaryRule) error {
	// 写入规则类型
	if err := binary.Write(writer, binary.LittleEndian, rule.ruleType); err != nil {
		return err
	}

	// 根据规则类型编码 payload
	switch rule.ruleType {
	case BinaryTypeDomain, BinaryTypeDomainSuffix:
		return w.writeDomainPayload(writer, rule.payload)
	case BinaryTypeIPCIDR, BinaryTypeIPCIDR6:
		return w.writeIPPayload(writer, rule.payload, rule.ruleType == BinaryTypeIPCIDR6)
	}

	return nil
}

// writeDomainPayload 写入域名 payload（带后缀压缩）
func (w *PrismBinaryWriter) writeDomainPayload(writer io.Writer, domain string) error {
	compressed, suffixID := w.compressDomain(domain)

	// 写入后缀 ID
	if err := binary.Write(writer, binary.LittleEndian, suffixID); err != nil {
		return err
	}

	// 写入压缩后的域名长度（使用 varint）
	if err := w.writeVarint(writer, uint64(len(compressed))); err != nil {
		return err
	}

	// 写入压缩后的域名
	_, err := writer.Write([]byte(compressed))
	return err
}

// writeIPPayload 写入 IP CIDR payload（二进制编码）
func (w *PrismBinaryWriter) writeIPPayload(writer io.Writer, cidr string, isIPv6 bool) error {
	if !w.ipBinaryEnabled {
		// 降级为字符串编码
		if err := w.writeVarint(writer, uint64(len(cidr))); err != nil {
			return err
		}
		_, err := writer.Write([]byte(cidr))
		return err
	}

	// 解析 CIDR
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		// 解析失败，降级为字符串编码
		if err := w.writeVarint(writer, uint64(len(cidr))); err != nil {
			return err
		}
		_, err := writer.Write([]byte(cidr))
		return err
	}

	// 写入标志位：0xFF 表示二进制编码
	if err := binary.Write(writer, binary.LittleEndian, uint8(0xFF)); err != nil {
		return err
	}

	// 写入 IP 地址（4 或 16 字节）
	addr := prefix.Addr()
	if _, err := writer.Write(addr.AsSlice()); err != nil {
		return err
	}

	// 写入前缀长度（1 字节）
	if err := binary.Write(writer, binary.LittleEndian, uint8(prefix.Bits())); err != nil {
		return err
	}

	return nil
}

// PrismBinaryReader Prism 二进制格式读取器
type PrismBinaryReader struct {
	gzipEnabled     bool
	varintEnabled   bool
	ipBinaryEnabled bool
	reverseSuffix   []string
}

// ReadPrismBinary 读取 Prism 二进制格式
func ReadPrismBinary(r io.Reader) (map[string][]Rule, error) {
	reader := &PrismBinaryReader{}

	// 读取并验证文件头
	flags, err := reader.readHeader(r)
	if err != nil {
		return nil, err
	}

	// 解析 flags
	reader.gzipEnabled = (flags & 0x01) != 0
	reader.varintEnabled = (flags & 0x02) != 0
	reader.ipBinaryEnabled = (flags & 0x04) != 0

	// 构建反向字典
	reader.reverseSuffix = domainSuffixDict

	// 如果启用 gzip，解压
	var dataReader io.Reader = r
	if reader.gzipEnabled {
		gzReader, err := gzip.NewReader(r)
		if err != nil {
			return nil, fmt.Errorf("create gzip reader: %w", err)
		}
		defer gzReader.Close()
		dataReader = gzReader
	}

	// 读取规则数据
	return reader.readRules(dataReader)
}

// readHeader 读取并验证文件头
func (r *PrismBinaryReader) readHeader(reader io.Reader) (flags uint8, err error) {
	// 读取 Magic Number
	magic := make([]byte, 5)
	if _, err := io.ReadFull(reader, magic); err != nil {
		return 0, fmt.Errorf("read magic: %w", err)
	}
	if !bytes.Equal(magic, []byte(MagicNumber)) {
		return 0, fmt.Errorf("invalid magic number: %s", magic)
	}

	// 读取 Version
	var version uint8
	if err := binary.Read(reader, binary.LittleEndian, &version); err != nil {
		return 0, fmt.Errorf("read version: %w", err)
	}
	if version != Version {
		return 0, fmt.Errorf("unsupported version: %d (expected %d)", version, Version)
	}

	// 读取 Flags
	if err := binary.Read(reader, binary.LittleEndian, &flags); err != nil {
		return 0, fmt.Errorf("read flags: %w", err)
	}

	// 跳过 Reserved
	reserved := make([]byte, 9)
	if _, err := io.ReadFull(reader, reserved); err != nil {
		return 0, fmt.Errorf("read reserved: %w", err)
	}

	return flags, nil
}

// readRules 读取规则数据
func (r *PrismBinaryReader) readRules(reader io.Reader) (map[string][]Rule, error) {
	// 读取 Action 数量
	var actionCount uint16
	if err := binary.Read(reader, binary.LittleEndian, &actionCount); err != nil {
		return nil, fmt.Errorf("read action count: %w", err)
	}

	rulesByAction := make(map[string][]Rule, actionCount)

	// 读取每个 Action
	for i := 0; i < int(actionCount); i++ {
		action, rules, err := r.readAction(reader)
		if err != nil {
			return nil, fmt.Errorf("read action %d: %w", i, err)
		}
		rulesByAction[action] = rules
	}

	return rulesByAction, nil
}

// readAction 读取单个 Action 及其规则
func (r *PrismBinaryReader) readAction(reader io.Reader) (string, []Rule, error) {
	// 读取 Action 名称长度
	var actionNameLen uint8
	if err := binary.Read(reader, binary.LittleEndian, &actionNameLen); err != nil {
		return "", nil, err
	}

	// 读取 Action 名称
	actionNameBytes := make([]byte, actionNameLen)
	if _, err := io.ReadFull(reader, actionNameBytes); err != nil {
		return "", nil, err
	}
	action := string(actionNameBytes)

	// 读取规则数量
	ruleCount, err := readVarint(reader, r.varintEnabled)
	if err != nil {
		return "", nil, err
	}

	// 读取规则列表
	rules := make([]Rule, 0, ruleCount)
	payload := prism.ParsePayload(action)

	for j := uint64(0); j < ruleCount; j++ {
		rule, err := r.readRule(reader, payload)
		if err != nil {
			return "", nil, fmt.Errorf("read rule %d: %w", j, err)
		}
		rules = append(rules, rule)
	}

	return action, rules, nil
}

// readRule 读取单条规则
func (r *PrismBinaryReader) readRule(reader io.Reader, payload prism.Payload) (Rule, error) {
	// 读取规则类型
	var ruleType BinaryRuleType
	if err := binary.Read(reader, binary.LittleEndian, &ruleType); err != nil {
		return nil, err
	}

	// 根据规则类型读取 payload
	switch ruleType {
	case BinaryTypeDomain:
		content, err := r.readDomainPayload(reader)
		if err != nil {
			return nil, err
		}
		return NewDomain(content, payload), nil

	case BinaryTypeDomainSuffix:
		content, err := r.readDomainPayload(reader)
		if err != nil {
			return nil, err
		}
		return NewDomainSuffix(content, payload), nil

	case BinaryTypeIPCIDR, BinaryTypeIPCIDR6:
		content, err := r.readIPPayload(reader, ruleType == BinaryTypeIPCIDR6)
		if err != nil {
			return nil, err
		}
		return NewIPCIDR(content, payload, false)

	default:
		return nil, fmt.Errorf("unknown rule type: %d", ruleType)
	}
}

// readDomainPayload 读取域名 payload
func (r *PrismBinaryReader) readDomainPayload(reader io.Reader) (string, error) {
	// 读取后缀 ID
	var suffixID uint8
	if err := binary.Read(reader, binary.LittleEndian, &suffixID); err != nil {
		return "", err
	}

	// 读取压缩域名长度
	compressedLen, err := readVarint(reader, r.varintEnabled)
	if err != nil {
		return "", err
	}

	// 读取压缩域名
	compressedBytes := make([]byte, compressedLen)
	if _, err := io.ReadFull(reader, compressedBytes); err != nil {
		return "", err
	}
	compressed := string(compressedBytes)

	// 解压域名
	if suffixID == 255 {
		return compressed, nil
	}
	if int(suffixID) < len(r.reverseSuffix) {
		return compressed + r.reverseSuffix[suffixID], nil
	}
	return compressed, nil
}

// readIPPayload 读取 IP CIDR payload
func (r *PrismBinaryReader) readIPPayload(reader io.Reader, isIPv6 bool) (string, error) {
	// 读取第一个字节（标志位或长度）
	var firstByte uint8
	if err := binary.Read(reader, binary.LittleEndian, &firstByte); err != nil {
		return "", err
	}

	// 如果是 0xFF，说明是二进制编码
	if firstByte == 0xFF && r.ipBinaryEnabled {
		// 读取 IP 地址
		ipSize := 4
		if isIPv6 {
			ipSize = 16
		}
		ipBytes := make([]byte, ipSize)
		if _, err := io.ReadFull(reader, ipBytes); err != nil {
			return "", err
		}

		addr, ok := netip.AddrFromSlice(ipBytes)
		if !ok {
			return "", fmt.Errorf("invalid IP address")
		}

		// 读取前缀长度
		var prefixLen uint8
		if err := binary.Read(reader, binary.LittleEndian, &prefixLen); err != nil {
			return "", err
		}

		prefix := netip.PrefixFrom(addr, int(prefixLen))
		return prefix.String(), nil
	}

	// 否则是字符串编码（firstByte 是 varint 的第一个字节）
	// 需要继续读取 varint
	length := uint64(firstByte)
	if r.varintEnabled && (firstByte&0x80) != 0 {
		// 多字节 varint，需要继续读取
		buf := []byte{firstByte}
		for {
			var b byte
			if err := binary.Read(reader, binary.LittleEndian, &b); err != nil {
				return "", err
			}
			buf = append(buf, b)
			if (b & 0x80) == 0 {
				break
			}
		}
		length, _ = binary.Uvarint(buf)
	}

	// 读取字符串
	cidrBytes := make([]byte, length)
	if _, err := io.ReadFull(reader, cidrBytes); err != nil {
		return "", err
	}

	return string(cidrBytes), nil
}
