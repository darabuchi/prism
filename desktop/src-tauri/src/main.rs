// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

use std::sync::{Arc, Mutex};
use std::thread;
use std::time::Duration;
use tauri::{
    CustomMenuItem, Manager, SystemTray, SystemTrayEvent, SystemTrayMenu, SystemTrayMenuItem,
    WindowEvent,
};
use tauri::api::notification::Notification;

// Tauri 应用状态
#[derive(Default)]
struct AppState {
    core_process: Arc<Mutex<Option<std::process::Child>>>,
}

// 启动 Prism Core 服务
#[tauri::command]
async fn start_core_service() -> Result<String, String> {
    println!("Starting Prism Core service...");
    
    // 这里应该启动实际的 core 服务
    // 为了演示，我们模拟启动过程
    tokio::time::sleep(Duration::from_secs(2)).await;
    
    Ok("Core service started successfully".to_string())
}

// 停止 Prism Core 服务
#[tauri::command]
async fn stop_core_service() -> Result<String, String> {
    println!("Stopping Prism Core service...");
    
    // 这里应该停止实际的 core 服务
    tokio::time::sleep(Duration::from_secs(1)).await;
    
    Ok("Core service stopped successfully".to_string())
}

// 获取系统信息
#[tauri::command]
async fn get_system_info() -> Result<serde_json::Value, String> {
    let info = serde_json::json!({
        "os": std::env::consts::OS,
        "arch": std::env::consts::ARCH,
        "version": env!("CARGO_PKG_VERSION"),
        "uptime": "2h 34m"
    });
    
    Ok(info)
}

// 显示系统通知
#[tauri::command]
async fn show_notification(title: String, body: String) -> Result<(), String> {
    Notification::new("com.prism.desktop")
        .title(&title)
        .body(&body)
        .show()
        .map_err(|e| e.to_string())?;
    
    Ok(())
}

// 检查更新
#[tauri::command]
async fn check_for_updates() -> Result<serde_json::Value, String> {
    // 模拟检查更新
    tokio::time::sleep(Duration::from_secs(1)).await;
    
    let update_info = serde_json::json!({
        "has_update": false,
        "current_version": env!("CARGO_PKG_VERSION"),
        "latest_version": env!("CARGO_PKG_VERSION")
    });
    
    Ok(update_info)
}

// 创建系统托盘
fn create_system_tray() -> SystemTray {
    let show = CustomMenuItem::new("show".to_string(), "显示主窗口");
    let hide = CustomMenuItem::new("hide".to_string(), "隐藏主窗口");
    let separator = SystemTrayMenuItem::Separator;
    let quit = CustomMenuItem::new("quit".to_string(), "退出应用");
    
    let tray_menu = SystemTrayMenu::new()
        .add_item(show)
        .add_item(hide)
        .add_native_item(separator)
        .add_item(quit);
    
    SystemTray::new().with_menu(tray_menu)
}

// 处理系统托盘事件
fn handle_system_tray_event(app: &tauri::AppHandle, event: SystemTrayEvent) {
    match event {
        SystemTrayEvent::LeftClick { .. } => {
            // 左键点击显示/隐藏主窗口
            if let Some(window) = app.get_window("main") {
                if window.is_visible().unwrap_or(false) {
                    let _ = window.hide();
                } else {
                    let _ = window.show();
                    let _ = window.set_focus();
                }
            }
        }
        SystemTrayEvent::MenuItemClick { id, .. } => {
            match id.as_str() {
                "show" => {
                    if let Some(window) = app.get_window("main") {
                        let _ = window.show();
                        let _ = window.set_focus();
                    }
                }
                "hide" => {
                    if let Some(window) = app.get_window("main") {
                        let _ = window.hide();
                    }
                }
                "quit" => {
                    app.exit(0);
                }
                _ => {}
            }
        }
        _ => {}
    }
}

fn main() {
    // 初始化日志
    env_logger::init();
    
    tauri::Builder::default()
        .manage(AppState::default())
        .system_tray(create_system_tray())
        .on_system_tray_event(handle_system_tray_event)
        .invoke_handler(tauri::generate_handler![
            start_core_service,
            stop_core_service,
            get_system_info,
            show_notification,
            check_for_updates
        ])
        .on_window_event(|event| {
            match event.event() {
                WindowEvent::CloseRequested { api, .. } => {
                    // 关闭窗口时隐藏到系统托盘而不是退出
                    event.window().hide().unwrap();
                    api.prevent_close();
                }
                _ => {}
            }
        })
        .setup(|app| {
            // 应用启动时的初始化工作
            let app_handle = app.handle();
            
            // 在后台线程中启动核心服务监控
            thread::spawn(move || {
                loop {
                    // 检查核心服务状态
                    thread::sleep(Duration::from_secs(30));
                    
                    // 这里可以添加核心服务的健康检查逻辑
                    // 如果服务异常，可以发送通知或自动重启
                }
            });
            
            println!("Prism Desktop application started");
            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}