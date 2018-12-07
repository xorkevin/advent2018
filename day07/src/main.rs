extern crate regex;

use regex::Regex;
use std::cmp;
use std::collections::{BinaryHeap, HashMap, HashSet};
use std::error;
use std::fmt;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";
const NUM_WORKERS: usize = 5;

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

#[derive(Debug)]
struct Task {
    id: char,
    deps: Vec<char>,
    next: Vec<char>,
}

impl Task {
    fn new(id: char) -> Task {
        Task {
            id: id,
            deps: Vec::new(),
            next: Vec::new(),
        }
    }

    fn add_dep(&mut self, id: char) {
        self.deps.push(id);
    }

    fn add_next(&mut self, id: char) {
        self.next.push(id);
    }

    fn can_start(&self, finished: &HashSet<char>) -> bool {
        for i in self.deps.iter() {
            if !finished.contains(&i) {
                return false;
            }
        }
        true
    }

    fn cost(&self) -> u32 {
        (self.id as u32) - ('A' as u32) + 61
    }
}

impl PartialEq for Task {
    fn eq(&self, other: &Task) -> bool {
        self.id == other.id
    }
}
impl Eq for Task {}

impl Ord for Task {
    fn cmp(&self, other: &Task) -> cmp::Ordering {
        other.id.cmp(&self.id)
    }
}

impl PartialOrd for Task {
    fn partial_cmp(&self, other: &Task) -> Option<cmp::Ordering> {
        Some(self.cmp(other))
    }
}

#[derive(Debug)]
struct Process<'a> {
    task: &'a Task,
    cost: u32,
}

impl<'a> Process<'a> {
    fn new(task: &'a Task) -> Process {
        Process {
            task: task,
            cost: task.cost(),
        }
    }
}

fn main() -> Result<(), Box<error::Error>> {
    let re = Regex::new(r"^Step (?P<dep>[A-Z]) .* (?P<next>[A-Z]) .*$")?;
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let (tasks, start) = {
        let mut tasks = HashMap::new();
        let mut nextset = HashSet::new();
        for line in reader.lines() {
            let l = line?;
            let caps = re
                .captures(&l)
                .ok_or(BasicError::new("Regex cannot capture line"))?;
            let dep = caps
                .name("dep")
                .ok_or(BasicError::new("Dep does not exist"))?
                .as_str()
                .chars()
                .next()
                .ok_or(BasicError::new("Failed to get dep character"))?;
            let next = caps
                .name("next")
                .ok_or(BasicError::new("Dep does not exist"))?
                .as_str()
                .chars()
                .next()
                .ok_or(BasicError::new("Failed to get next character"))?;

            nextset.insert(next);
            {
                let k = tasks.entry(dep).or_insert(Task::new(dep));
                k.add_next(next);
            }
            {
                let k = tasks.entry(next).or_insert(Task::new(next));
                k.add_dep(dep);
            }
        }

        let mut start = tasks
            .keys()
            .filter(|x| !nextset.contains(x))
            .map(|&x| x)
            .collect::<Vec<_>>();
        start.sort();
        (tasks, start)
    };

    {
        let mut order = Vec::new();
        let mut openlist = BinaryHeap::new();
        for i in start.iter() {
            openlist.push(tasks.get(i).ok_or("Failed to get task")?);
        }
        let mut closedlist = HashSet::new();
        while let Some(i) = openlist.pop() {
            order.push(i.id);
            closedlist.insert(i.id);
            for taskid in i.next.iter() {
                if closedlist.contains(&taskid) {
                    continue;
                }
                let task = tasks.get(&taskid).ok_or("Failed to get task")?;
                if !task.can_start(&closedlist) {
                    continue;
                }
                openlist.push(task);
            }
        }
        println!("{}", order.iter().collect::<String>());
    }

    {
        let mut elapsed = 0;

        let mut openlist = BinaryHeap::new();
        for i in start.iter() {
            openlist.push(tasks.get(i).ok_or("Failed to get task")?);
        }
        let mut closedlist = HashSet::new();
        let mut current_work = Vec::new();

        'workloop: loop {
            while current_work.len() < NUM_WORKERS {
                if let Some(i) = openlist.pop() {
                    current_work.push(Process::new(i));
                } else {
                    if current_work.len() == 0 {
                        break 'workloop;
                    }
                    break;
                }
            }
            let min_time = current_work
                .iter()
                .min_by_key(|x| x.cost)
                .ok_or("Failed to find minimum cost")?
                .cost;
            elapsed += min_time;
            let mut next_work_ids = HashSet::new();
            let mut next_work = Vec::new();
            let mut done_work = Vec::new();
            while let Some(mut i) = current_work.pop() {
                i.cost -= min_time;
                if i.cost == 0 {
                    closedlist.insert(i.task.id);
                    done_work.push(i);
                } else {
                    next_work_ids.insert(i.task.id);
                    next_work.push(i);
                }
            }
            current_work = next_work;

            for i in done_work.iter() {
                for taskid in i.task.next.iter() {
                    if next_work_ids.contains(&taskid) {
                        continue;
                    }
                    if closedlist.contains(&taskid) {
                        continue;
                    }
                    let task = tasks.get(&taskid).ok_or("Failed to get task")?;
                    if !task.can_start(&closedlist) {
                        continue;
                    }
                    openlist.push(task);
                }
            }
        }
        println!("{}", elapsed);
    }

    Ok(())
}
