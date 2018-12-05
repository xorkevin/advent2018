use std::error;
use std::fmt;
use std::fs::File;
use std::io::prelude::*;

const PUZZLEINPUT: &str = "input.txt";

#[derive(Debug)]
struct BasicError {
    msg: String,
}

impl BasicError {
    fn new(msg: &str) -> BasicError {
        BasicError {
            msg: msg.to_string(),
        }
    }
}

impl fmt::Display for BasicError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", self.msg)
    }
}

impl error::Error for BasicError {
    fn description(&self) -> &str {
        &self.msg
    }

    fn cause(&self) -> Option<&error::Error> {
        None
    }
}

fn main() -> Result<(), Box<error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let init_state = file
        .bytes()
        .filter_map(|x| x.ok())
        .map(|x| x as char)
        .collect::<Vec<_>>();

    Ok(())
}
