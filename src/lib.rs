extern crate libc;

use libc::c_char;
use std::ffi::CStr;
use std::ffi::CString;
use std::process::Command;

#[derive(Debug, Clone)]
#[repr(C)]
pub struct Response {
    pub output: *const c_char,
    pub err: *const c_char,
}

impl Response {
    pub fn boxed<T: Into<Vec<u8>>>(output: T, err: T) -> *mut Response {
        let output_ptr = string_to_c_char(output);
        let err_ptr = string_to_c_char(err);

        let resp = Response {
            output: output_ptr,
            err: err_ptr,
        };
        return Box::into_raw(Box::new(resp));
    }
}

fn string_to_c_char<T: Into<Vec<u8>>>(value: T) -> *const c_char {
    let value_str = CString::new(value).unwrap();
    let value_ptr = value_str.as_ptr();
    std::mem::forget(value_str);
    return value_ptr;
}

#[no_mangle]
pub extern "C" fn getCommandOutput(command: *const c_char) -> *mut Response {
    let buf_command = unsafe { CStr::from_ptr(command).to_bytes() };
    let str_command = String::from_utf8(buf_command.to_vec()).unwrap();
    let cmd = Command::new(str_command)
        .arg("--version")
        .output();
    match cmd {
        Ok(cmd) => {
            return Response::boxed(cmd.stdout, cmd.stderr);
        },
        Err(_e) => {
            return Response::boxed("", "Command::new failed");
        }
    }
}

#[no_mangle]
pub unsafe extern "C" fn DestroyResponse(resp: *mut Response) {
    if !resp.is_null() {
        drop(Box::from_raw(resp));
    }
}
