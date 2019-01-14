use regex::Regex;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";
const PUZZLEINPUT2: &str = "input2.txt";

enum OpCode {
    Addr,
    Addi,
    Mulr,
    Muli,
    Banr,
    Bani,
    Borr,
    Bori,
    Setr,
    Seti,
    Gtir,
    Gtri,
    Gtrr,
    Eqir,
    Eqri,
    Eqrr,
}

impl OpCode {
    fn from_isize(opcode: isize) -> Option<OpCode> {
        match opcode {
            0 => Some(OpCode::Addr),
            1 => Some(OpCode::Addi),
            2 => Some(OpCode::Mulr),
            3 => Some(OpCode::Muli),
            4 => Some(OpCode::Banr),
            5 => Some(OpCode::Bani),
            6 => Some(OpCode::Borr),
            7 => Some(OpCode::Bori),
            8 => Some(OpCode::Setr),
            9 => Some(OpCode::Seti),
            10 => Some(OpCode::Gtir),
            11 => Some(OpCode::Gtri),
            12 => Some(OpCode::Gtrr),
            13 => Some(OpCode::Eqir),
            14 => Some(OpCode::Eqri),
            15 => Some(OpCode::Eqrr),
            _ => None,
        }
    }
}

fn exec_instr(opcode: isize, ai: isize, bi: isize, ci: isize, s: &mut Vec<isize>) {
    let a = ai as usize;
    let b = bi as usize;
    let c = ci as usize;
    s[c] = match OpCode::from_isize(opcode).expect("Failed to parse opcode") {
        OpCode::Addr => s[a] + s[b],
        OpCode::Addi => s[a] + bi,
        OpCode::Mulr => s[a] * s[b],
        OpCode::Muli => s[a] * bi,
        OpCode::Banr => s[a] & s[b],
        OpCode::Bani => s[a] & bi,
        OpCode::Borr => s[a] | s[b],
        OpCode::Bori => s[a] | bi,
        OpCode::Setr => s[a],
        OpCode::Seti => ai,
        OpCode::Gtir => {
            if ai > s[b] {
                1
            } else {
                0
            }
        }
        OpCode::Gtri => {
            if s[a] > bi {
                1
            } else {
                0
            }
        }
        OpCode::Gtrr => {
            if s[a] > s[b] {
                1
            } else {
                0
            }
        }
        OpCode::Eqir => {
            if ai == s[b] {
                1
            } else {
                0
            }
        }
        OpCode::Eqri => {
            if s[a] == bi {
                1
            } else {
                0
            }
        }
        OpCode::Eqrr => {
            if s[a] == s[b] {
                1
            } else {
                0
            }
        }
    };
}

struct TranslationTable {
    table: Vec<Vec<isize>>,
}

struct TestCase(Vec<isize>, (isize, isize, isize, isize), Vec<isize>);

impl TranslationTable {
    fn new() -> Self {
        Self {
            table: vec![vec![0; 16]; 16],
        }
    }

    fn provide(&mut self, TestCase(before, op, after): TestCase) -> usize {
        let mut count = 0;
        for i in 0..self.table.len() {
            let mut test = before.clone();
            exec_instr(i as isize, op.1, op.2, op.3, &mut test);
            if test != after {
                self.table[op.0 as usize][i] = 1;
                count += 1;
            }
        }
        count
    }

    fn gen_translation(&mut self) -> Vec<isize> {
        let mut translation = vec![16; 16];
        let mut k = 0;
        while k < 16 {
            for op in 0..self.table.len() {
                if translation[op] < 16 {
                    continue;
                }
                let opts = &self.table[op];
                let mut count = 0;
                let mut last = 0;
                for (n, &i) in opts.iter().enumerate() {
                    if i != 1 {
                        count += 1;
                        last = n;
                    }
                }
                if count == 1 {
                    for i in 0..self.table.len() {
                        self.table[i][last] = 1;
                    }
                    translation[op] = last as isize;
                    k += 1;
                    break;
                }
            }
        }
        translation
    }
}

fn main() {
    part2(part1());
}

fn parse_line(re: &Regex, line: &str) -> (isize, isize, isize, isize) {
    let caps = re.captures(&line).expect("Regex cannot capture line");
    let a = caps
        .name("a")
        .expect("a does not exist")
        .as_str()
        .parse::<isize>()
        .expect("Failed to parse a");
    let b = caps
        .name("b")
        .expect("b does not exist")
        .as_str()
        .parse::<isize>()
        .expect("Failed to parse b");
    let c = caps
        .name("c")
        .expect("c does not exist")
        .as_str()
        .parse::<isize>()
        .expect("Failed to parse c");
    let d = caps
        .name("d")
        .expect("d does not exist")
        .as_str()
        .parse::<isize>()
        .expect("Failed to parse d");
    (a, b, c, d)
}

fn part1() -> Vec<isize> {
    let re_before = Regex::new(r"^Before: \[(?P<a>\d+), (?P<b>\d+), (?P<c>\d+), (?P<d>\d+)\]$")
        .expect("Invalid regex");
    let re_after = Regex::new(r"^After:  ?\[(?P<a>\d+), (?P<b>\d+), (?P<c>\d+), (?P<d>\d+)\]$")
        .expect("Invalid regex");
    let re_op =
        Regex::new(r"^(?P<a>\d+) (?P<b>\d+) (?P<c>\d+) (?P<d>\d+)$").expect("Invalid regex");

    let file = File::open(PUZZLEINPUT).expect("Failed to open file");
    let reader = BufReader::new(file);
    let mut lines = reader.lines();

    let mut table = TranslationTable::new();
    let mut part1 = 0;

    while let Some(line) = lines.next() {
        let before = parse_line(&re_before, &line.expect("Failed to read line"));
        let op = parse_line(
            &re_op,
            &lines
                .next()
                .expect("Failed to read line")
                .expect("Failed to read line"),
        );
        let after = parse_line(
            &re_after,
            &lines
                .next()
                .expect("Failed to read line")
                .expect("Failed to read line"),
        );
        lines.next();
        if 16
            - table.provide(TestCase(
                vec![before.0, before.1, before.2, before.3],
                op,
                vec![after.0, after.1, after.2, after.3],
            ))
            > 2
        {
            part1 += 1;
        }
    }

    println!("{}", part1);

    table.gen_translation()
}

fn part2(translation: Vec<isize>) {
    let re_op =
        Regex::new(r"^(?P<a>\d+) (?P<b>\d+) (?P<c>\d+) (?P<d>\d+)$").expect("Invalid regex");
    let file = File::open(PUZZLEINPUT2).expect("Failed to open file");
    let reader = BufReader::new(file);

    let mut state = vec![0; 4];
    for line in reader.lines() {
        let op = parse_line(&re_op, &line.expect("Failed to read line"));
        exec_instr(translation[op.0 as usize], op.1, op.2, op.3, &mut state);
    }
    println!("{}", state[0]);
}
