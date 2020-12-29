extern crate libc;

use libc::c_char;
use std::env;
use std::ffi::CStr;
use std::ffi::CString;
use std::fs;
use std::process::Command;

#[derive(Debug, Clone)]
#[repr(C)]
pub struct Response {
    pub output: *const c_char,
    pub err: *const c_char,
}

impl Response {
    pub fn new<T: Into<Vec<u8>>>(output: T, err: T) -> Response {
        let output_ptr = string_to_c_char(output);
        let err_ptr = string_to_c_char(err);

        Response {
            output: output_ptr,
            err: err_ptr,
        }
    }
}

fn string_to_c_char<T: Into<Vec<u8>>>(value: T) -> *const c_char {
    let value_str = CString::new(value).unwrap();
    let value_ptr = value_str.as_ptr();
    std::mem::forget(value_str);
    return value_ptr;
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

#[no_mangle]
pub extern "C" fn getCommandOutput(command: *const c_char) -> *mut Response {
    let buf_command = unsafe { CStr::from_ptr(command).to_bytes() };
    let str_command = String::from_utf8(buf_command.to_vec()).unwrap();
    if !is_program_in_path(&str_command) {
        let resp = Response::new("", "program not in path");
        return Box::into_raw(Box::new(resp))
    }
    let output = Command::new(str_command)
        .arg("Hello World")
        .output()
        .expect("Command::new failed");
    let resp = Response::new(output.stdout, output.stderr);
    return Box::into_raw(Box::new(resp))
}

#[no_mangle]
pub unsafe extern "C" fn DestroyResponse(resp: *mut Response) {
    if !resp.is_null() {
        drop(Box::from_raw(resp));
    }
}
