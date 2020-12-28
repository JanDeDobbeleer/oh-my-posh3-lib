extern crate libc;

use libc::c_char;
use std::env;
use std::ffi::CStr;
use std::ffi::CString;
use std::fs;
use std::process::Command;

#[repr(C)]
pub struct Response {
    output: *const c_char,
    err: *const c_char,
}

fn is_program_in_path(program: &str) -> bool {
    if let Ok(path) = env::var("PATH") {
        for p in path.split(":") {
            let p_str = format!("{}/{}", p, program);
            if fs::metadata(p_str).is_ok() {
                return true;
            }
        }
    }
    false
}

fn string_to_c_char<T: Into<Vec<u8>>>(value: T) -> *const c_char {
    let value_str = CString::new(value).unwrap();
    let value_ptr = value_str.as_ptr();
    std::mem::forget(value_str);
    return value_ptr;
}

#[no_mangle]
pub extern "C" fn getCommandOutput(command: *const c_char) -> *const c_char {
    let buf_command = unsafe { CStr::from_ptr(command).to_bytes() };
    let str_command = String::from_utf8(buf_command.to_vec()).unwrap();
    if !is_program_in_path(&str_command) {
        return string_to_c_char("err: program not in path")
    }
    let output = Command::new(str_command)
        .arg("Hello World")
        .output()
        .expect("Command::new failed");
    if !output.status.success() {
        let error = String::from_utf8(output.stderr).unwrap();
        return string_to_c_char(format!("err: {}", error));
    }
    return string_to_c_char(output.stdout);
}
