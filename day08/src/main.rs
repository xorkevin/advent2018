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

struct Node {
    children: Vec<Node>,
    metadata: Vec<usize>,
}

impl Node {
    fn new() -> Node {
        Node {
            children: Vec::new(),
            metadata: Vec::new(),
        }
    }

    fn add_child(&mut self, child: Node) {
        self.children.push(child);
    }

    fn add_metadata(&mut self, data: usize) {
        self.metadata.push(data);
    }

    fn sum(&self) -> usize {
        if self.children.len() == 0 {
            self.metadata.iter().sum()
        } else {
            self.metadata
                .iter()
                .filter_map(|&x| self.children.get(x - 1))
                .map(|x| x.sum())
                .sum()
        }
    }
}

fn main() -> Result<(), Box<error::Error>> {
    let nums = {
        let mut file = File::open(PUZZLEINPUT)?;
        let mut buffer = String::new();
        file.read_to_string(&mut buffer)?;
        let mut nums = Vec::new();
        for i in buffer.split(" ") {
            nums.push(i.trim().parse::<usize>()?);
        }
        nums
    };

    let mut iter = nums.iter();

    let (n, sum) = process_stream(&mut iter)?;
    println!("{}", sum);
    println!("{}", n.sum());

    Ok(())
}

fn process_stream(nums: &mut Iterator<Item = &usize>) -> Result<(Node, usize), BasicError> {
    let &num_children = nums
        .next()
        .ok_or(BasicError::new("Failed to get num children"))?;
    let &num_metadata = nums
        .next()
        .ok_or(BasicError::new("Failed to get num metadata"))?;

    let mut sum = 0;
    let mut node = Node::new();

    for _ in 0..num_children {
        let (child, s) = process_stream(nums)?;
        node.add_child(child);
        sum += s;
    }
    for _ in 0..num_metadata {
        let &data = nums
            .next()
            .ok_or(BasicError::new("Failed to get metadata"))?;
        node.add_metadata(data);
        sum += data;
    }
    Ok((node, sum))
}
