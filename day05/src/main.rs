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

fn is_reactive(a: &char, b: &char) -> bool {
    a != b && a.to_ascii_uppercase() == b.to_ascii_uppercase()
}

fn reduce(state: &Vec<char>) -> Vec<char> {
    state.iter().fold(Vec::new(), |mut acc, i| {
        if {
            match acc.last() {
                Some(x) => is_reactive(i, x),
                None => false,
            }
        } {
            acc.pop();
        } else {
            acc.push(*i);
        }
        acc
    })
}

fn is_not_char(a: &char, b: &char) -> bool {
    a != b && a.to_ascii_uppercase() != b.to_ascii_uppercase()
}

fn main() -> Result<(), Box<error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let init_state = file
        .bytes()
        .filter_map(|x| x.ok())
        .map(|x| x as char)
        .filter(char::is_ascii_alphanumeric)
        .collect::<Vec<_>>();

    println!("{}", reduce(&init_state).len());

    let mut max = 99999999;
    for i in ('A' as u8)..=('Z' as u8) {
        let c = i as char;
        let l = reduce(
            &init_state
                .iter()
                .filter(|x| is_not_char(x, &c))
                .map(|&x| x)
                .collect::<Vec<_>>(),
        ).len();
        if l < max {
            max = l;
        }
    }
    println!("{}", max);

    Ok(())
}
