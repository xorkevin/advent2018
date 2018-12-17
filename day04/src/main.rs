use std::collections::HashMap;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() {
    let re = regex::Regex::new(r"^\[\d+-\d+-\d+ \d+:(?P<minute>\d+)\] (?P<action>[A-Za-z0-9# ]+)$")
        .expect("Invalid regex");
    let file = File::open(PUZZLEINPUT).expect("Failed to open file");
    let reader = BufReader::new(file);

    let timesheet = {
        let mut timesheet = HashMap::new();
        let mut current_guard = 0;
        let mut asleep_min = 0;
        for line in reader.lines() {
            let l = line.expect("Failed to read line");
            let caps = re.captures(&l).expect("Invalid input 1");
            let minute = caps
                .name("minute")
                .expect("Invalid input 2")
                .as_str()
                .parse::<usize>()
                .expect("Failed to parse minute");
            let action = caps
                .name("action")
                .expect("Invalid input 3")
                .as_str()
                .split_whitespace()
                .collect::<Vec<_>>();
            match action[0] {
                "Guard" => {
                    current_guard = action[1][1..]
                        .parse::<usize>()
                        .expect("Failed to parse guard")
                }
                "falls" => asleep_min = minute,
                "wakes" => {
                    let guard = timesheet.entry(current_guard).or_insert((0, vec![0; 60]));
                    let (ref mut total, ref mut times) = *guard;
                    *total += minute - asleep_min;
                    for i in asleep_min..minute {
                        times[i] += 1;
                    }
                }
                _ => panic!("Invalid input 4"),
            }
        }
        timesheet
    };

    {
        let (guardid, (_, guardtimes)) = timesheet
            .iter()
            .max_by_key(|(_, (x, _))| x)
            .expect("Failed to find max");
        let (minute, _) = guardtimes
            .iter()
            .enumerate()
            .max_by_key(|(_, &x)| x)
            .expect("Failed to find max");
        println!("{}", guardid * minute);
    };

    {
        let (guardid, (minute, _)) = timesheet
            .iter()
            .map(|(guardid, (_, times))| {
                (
                    guardid,
                    times
                        .iter()
                        .enumerate()
                        .max_by_key(|(_, &x)| x)
                        .expect("Failed to find max"),
                )
            })
            .max_by_key(|(_, (_, &x))| x)
            .expect("Failed to find max");
        println!("{}", guardid * minute);
    };
}
