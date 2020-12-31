extern crate libc;

use libc::c_char;
use libc::c_int;
use std::ffi::CStr;
use std::ffi::CString;
use std::process::Command;
use std::process::ExitStatus;
use std::ptr;

#[derive(Debug, Clone)]
#[repr(C)]
pub struct Response {
    pub output: *const c_char,
    pub err: *const c_char,
    pub status_code: c_int,
}

impl Response {
    pub fn boxed<T: Into<Vec<u8>>>(output: T, err: T, status_code: i32) -> *mut Response {
        let output_ptr = string_to_c_char(output);
        let err_ptr = string_to_c_char(err);
        let resp = Response {
            output: output_ptr,
            err: err_ptr,
            status_code: status_code as c_int,
        };
        return Box::into_raw(Box::new(resp));
    }
}

fn c_char_to_string(char_c: *const c_char) -> String {
    if char_c == ptr::null() {
        return "".to_string();
    }
    let buf = unsafe { CStr::from_ptr(char_c).to_bytes() };
    return String::from_utf8(buf.to_vec()).unwrap();
}

fn string_to_c_char<T: Into<Vec<u8>>>(value: T) -> *const c_char {
    let value_str = CString::new(value).unwrap();
    let value_ptr = value_str.as_ptr();
    std::mem::forget(value_str);
    return value_ptr;
}

fn get_safe_status_code(status: ExitStatus) -> i32 {
    match status.code() {
        Some(code) => return code,
        None       => return 0
    }
}

#[no_mangle]
pub extern "C" fn getCommandOutput(command: *const c_char, args: *const c_char) -> *mut Response {
    let str_command = c_char_to_string(command);
    let mut cmd = Command::new(str_command);
    if args != ptr::null() {
        let args_str = c_char_to_string(args);
        let args_split = args_str.split(";");
        cmd.args(args_split);
    }
    let output = cmd.output();
    match output {
        Ok(cmd) => {
            let status_code = get_safe_status_code(cmd.status);
            return Response::boxed(cmd.stdout, cmd.stderr, status_code);
        }
        Err(_e) => {
            return Response::boxed("", "Command::new failed", 0);
        }
    }
}

#[no_mangle]
pub unsafe extern "C" fn DestroyResponse(resp: *mut Response) {
    if !resp.is_null() {
        drop(Box::from_raw(resp));
    }
}
