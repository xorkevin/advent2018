use std::collections::VecDeque;
use std::error::Error;
use std::fmt;

const PLAYER_COUNT: usize = 486;
const PART1_VALUE: usize = 70833;
const LAST_VALUE: usize = PART1_VALUE * 100;

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

impl Error for BasicError {
    fn description(&self) -> &str {
        &self.msg
    }

    fn cause(&self) -> Option<&Error> {
        None
    }
}

fn main() -> Result<(), Box<Error>> {
    let mut score = vec![0; PLAYER_COUNT];
    let mut current = VecDeque::with_capacity(LAST_VALUE);
    current.push_back(0);
    for i in 0..LAST_VALUE {
        let player = i % PLAYER_COUNT;
        let value = i + 1;
        if value % 23 == 0 {
            for _ in 0..6 {
                let k = current.pop_back().ok_or(BasicError::new(
                    "Failed to get node to rotate counterclockwise",
                ))?;
                current.push_front(k);
            }
            let k = current
                .pop_back()
                .ok_or(BasicError::new("Failed to get node to remove"))?;
            score[player] += value + k;
        } else {
            for _ in 0..2 {
                let k = current
                    .pop_front()
                    .ok_or(BasicError::new("Failed to get node to rotate clockwise"))?;
                current.push_back(k);
            }
            current.push_front(value);
        }
    }

    println!(
        "{}",
        score
            .iter()
            .max()
            .ok_or(BasicError::new("Failed to get max"))?
    );
    Ok(())
}
