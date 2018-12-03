extern crate regex;

use std::error;
use std::fmt;
use std::fs;
use std::io;
use std::io::prelude::*;

const PUZZLEINPUT: &str = "input.txt";
const WIDTH: usize = 1000;

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
    let file = fs::File::open(PUZZLEINPUT)?;
    let reader = io::BufReader::new(file);

    let re = regex::Regex::new(
        r"^#(?P<claimid>\d+) @ (?P<col>\d+),(?P<row>\d+): (?P<width>\d+)x(?P<height>\d+)$",
    ).unwrap();
    let mut claims_count = vec![0; WIDTH * WIDTH];
    let mut claims = vec![0; WIDTH * WIDTH];
    let mut banned_claims = vec![false; 1253];

    for line in reader.lines() {
        let s = line?;
        let caps = re.captures(&s).ok_or(BasicError::new("Invalid input"))?;
        let claimid = caps
            .name("claimid")
            .ok_or(BasicError::new("Invalid input"))?
            .as_str()
            .parse::<usize>()?;
        let col = caps
            .name("col")
            .ok_or(BasicError::new("Invalid input"))?
            .as_str()
            .parse::<usize>()?;
        let row = caps
            .name("row")
            .ok_or(BasicError::new("Invalid input"))?
            .as_str()
            .parse::<usize>()?;
        let width = caps
            .name("width")
            .ok_or(BasicError::new("Invalid input"))?
            .as_str()
            .parse::<usize>()?;
        let height = caps
            .name("height")
            .ok_or(BasicError::new("Invalid input"))?
            .as_str()
            .parse::<usize>()?;
        for i in row..(row + height) {
            for j in col..(col + width) {
                claims_count[i * WIDTH + j] += 1;
                let mut c2 = &mut claims[i * WIDTH + j];
                if *c2 > 0 {
                    banned_claims[claimid - 1] = true;
                    banned_claims[*c2 - 1] = true;
                }
                *c2 = claimid;
            }
        }
    }

    println!("{}", claims_count.iter().filter(|&i| *i > 1).count());

    for i in banned_claims
        .iter()
        .enumerate()
        .filter(|(_, &i)| !i)
        .map(|(n, _)| n)
    {
        println!("{}", i + 1);
    }

    Ok(())
}
