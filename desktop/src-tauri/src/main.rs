// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

use tauri::{
    AppHandle, CustomMenuItem, Manager, SystemTray, SystemTrayEvent, SystemTrayMenu,
    SystemTrayMenuItem, WindowEvent,
};

// 命令定义
#[tauri::command]
async fn get_system_info() -> Result<serde_json::Value, String> {
    // 获取系统信息
    let info = serde_json::json!({
        "platform": std::env::consts::OS,
        "arch": std::env::consts::ARCH,
        "version": env!("CARGO_PKG_VERSION")
    });
    Ok(info)
}

#[tauri::command]
async fn check_core_connection(url: String) -> Result<bool, String> {
    match reqwest::get(&format!("{}/api/v1/health", url)).await {
        Ok(response) => Ok(response.status().is_success()),
        Err(_) => Ok(false),
    }
}

#[tauri::command]
async fn minimize_to_tray(window: tauri::Window) -> Result<(), String> {
    window.hide().map_err(|e| e.to_string())?;
    Ok(())
}

#[tauri::command]
async fn show_from_tray(app_handle: AppHandle) -> Result<(), String> {
    if let Some(window) = app_handle.get_window("main") {
        window.show().map_err(|e| e.to_string())?;
        window.set_focus().map_err(|e| e.to_string())?;
    }
    Ok(())
}

fn create_system_tray() -> SystemTray {
    let quit = CustomMenuItem::new("quit".to_string(), "退出");
    let show = CustomMenuItem::new("show".to_string(), "显示");
    let hide = CustomMenuItem::new("hide".to_string(), "隐藏");
    
    let tray_menu = SystemTrayMenu::new()
        .add_item(show)
        .add_item(hide)
        .add_native_item(SystemTrayMenuItem::Separator)
        .add_item(quit);
    
    SystemTray::new().with_menu(tray_menu)
}

fn handle_system_tray_event(app: &AppHandle, event: SystemTrayEvent) {
    match event {
        SystemTrayEvent::LeftClick {
            position: _,
            size: _,
            ..
        } => {
            if let Some(window) = app.get_window("main") {
                if window.is_visible().unwrap_or(false) {
                    let _ = window.hide();
                } else {
                    let _ = window.show();
                    let _ = window.set_focus();
                }
            }
        }
        SystemTrayEvent::MenuItemClick { id, .. } => match id.as_str() {
            "quit" => {
                std::process::exit(0);
            }
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
            _ => {}
        },
        _ => {}
    }
}

fn handle_window_event(event: &WindowEvent) {
    match event {
        WindowEvent::CloseRequested { api, .. } => {
            // 阻止默认关闭行为，改为隐藏到托盘
            api.prevent_close();
            if let Some(window) = api.window().get_window("main") {
                let _ = window.hide();
            }
        }
        _ => {}
    }
}

fn main() {
    tauri::Builder::default()
        .system_tray(create_system_tray())
        .on_system_tray_event(handle_system_tray_event)
        .on_window_event(handle_window_event)
        .invoke_handler(tauri::generate_handler![
            get_system_info,
            check_core_connection,
            minimize_to_tray,
            show_from_tray
        ])
        .setup(|app| {
            // 应用启动时显示主窗口
            if let Some(window) = app.get_window("main") {
                let _ = window.show();
            }
            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}