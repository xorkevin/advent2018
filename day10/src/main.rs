extern crate regex;

use regex::Regex;
use std::error::Error;
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

impl Error for BasicError {
    fn description(&self) -> &str {
        &self.msg
    }

    fn cause(&self) -> Option<&Error> {
        None
    }
}

struct PointList {
    points: Vec<(isize, isize, isize, isize)>,
}

impl PointList {
    fn new() -> Self {
        Self { points: Vec::new() }
    }

    fn add(&mut self, point: (isize, isize, isize, isize)) {
        self.points.push(point);
    }

    fn step(&mut self, t: isize) {
        for &mut (ref mut x, ref mut y, vx, vy) in self.points.iter_mut() {
            *x += t * vx;
            *y += t * vy;
        }
    }

    fn size(&self) -> (usize, usize, isize, isize) {
        let (mut minx, mut miny, _, _) = self.points[0];
        let (mut maxx, mut maxy, _, _) = self.points[0];
        for &(x, y, _, _) in self.points.iter().skip(1) {
            if x < minx {
                minx = x;
            }
            if x > maxx {
                maxx = x;
            }
            if y < miny {
                miny = y;
            }
            if y > maxy {
                maxy = y;
            }
        }
        (
            (maxx - minx + 1) as usize,
            (maxy - miny + 1) as usize,
            minx,
            miny,
        )
    }

    fn print(&self) {
        let (sizex, sizey, minx, miny) = self.size();
        let mut board = vec![vec!['.'; sizex]; sizey];
        for &(x, y, _, _) in self.points.iter() {
            let x = (x - minx) as usize;
            let y = (y - miny) as usize;
            board[y][x] = '#';
        }
        for i in board {
            println!("{}", i.iter().collect::<String>());
        }
        println!("{} {} {} {}", sizex, sizey, minx, miny);
    }
}

fn main() -> Result<(), Box<Error>> {
    let re = Regex::new(
        r"^.*< ?(?P<posx>-?\d+),  ?(?P<posy>-?\d+)>.*< ?(?P<velx>-?\d+),  ?(?P<vely>-?\d+)>$",
    )?;
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut points = PointList::new();
    for line in reader.lines() {
        let l = line?;
        let caps = re
            .captures(&l)
            .ok_or(BasicError::new("Regex cannot capture line"))?;
        let posx = caps
            .name("posx")
            .ok_or(BasicError::new("Posx does not exist"))?
            .as_str()
            .parse::<isize>()?;
        let posy = caps
            .name("posy")
            .ok_or(BasicError::new("Posy does not exist"))?
            .as_str()
            .parse::<isize>()?;
        let velx = caps
            .name("velx")
            .ok_or(BasicError::new("Velx does not exist"))?
            .as_str()
            .parse::<isize>()?;
        let vely = caps
            .name("vely")
            .ok_or(BasicError::new("Vely does not exist"))?
            .as_str()
            .parse::<isize>()?;
        points.add((posx, posy, velx, vely));
    }

    let step = 10656;
    points.step(step);
    points.print();
    println!("Step: {}", step);

    Ok(())
}
