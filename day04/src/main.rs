use std::collections::HashMap;
use std::error;
use std::fmt;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

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
    let re = regex::Regex::new(r"^\[\d+-\d+-\d+ \d+:(?P<minute>\d+)\] (?P<action>[A-Za-z0-9# ]+)$")
        .unwrap();
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let timesheet = {
        let mut timesheet = HashMap::new();
        let mut current_guard = 0;
        let mut asleep_min = 0;
        for line in reader.lines() {
            let l = line?;
            let caps = re.captures(&l).ok_or(BasicError::new("Invalid input 1"))?;
            let minute = caps
                .name("minute")
                .ok_or(BasicError::new("Invalid input 2"))?
                .as_str()
                .parse::<usize>()?;
            let action = caps
                .name("action")
                .ok_or(BasicError::new("Invalid input 3"))?
                .as_str()
                .split_whitespace()
                .collect::<Vec<_>>();
            match action[0] {
                "Guard" => current_guard = action[1][1..].parse::<usize>()?,
                "falls" => asleep_min = minute,
                "wakes" => {
                    let guard = timesheet.entry(current_guard).or_insert((0, vec![0; 60]));
                    let (ref mut total, ref mut times) = *guard;
                    *total += minute - asleep_min;
                    for i in asleep_min..minute {
                        times[i] += 1;
                    }
                }
                _ => Err(BasicError::new("Invalid input 4"))?,
            }
        }
        timesheet
    };

    {
        let (guardid, (_, guardtimes)) = timesheet
            .iter()
            .max_by_key(|(_, (x, _))| x)
            .ok_or(BasicError::new("Failed to find max"))?;
        let (minute, _) = guardtimes
            .iter()
            .enumerate()
            .max_by_key(|(_, &x)| x)
            .ok_or(BasicError::new("Failed to find max"))?;
        println!("{}", guardid * minute);
    };

    {
        let (guardid, (minute, _)) = timesheet
            .iter()
            .filter_map(|(guardid, (_, times))| {
                match times.iter().enumerate().max_by_key(|(_, &x)| x) {
                    Some(x) => Some((guardid, x)),
                    _ => None,
                }
            }).max_by_key(|(_, (_, &x))| x)
            .ok_or(BasicError::new("Failed to find max"))?;
        println!("{}", guardid * minute);
    };

    Ok(())
}
