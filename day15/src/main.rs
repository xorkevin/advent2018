use std::cmp::Ordering;
use std::collections::{BinaryHeap, HashMap, HashSet};
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn distance(a: usize, b: usize) -> usize {
    if a > b {
        a - b
    } else {
        b - a
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
struct Pos(usize, usize);

impl Pos {
    fn up(&self) -> Self {
        Pos(self.0, self.1 - 1)
    }
    fn down(&self) -> Self {
        Pos(self.0, self.1 + 1)
    }
    fn left(&self) -> Self {
        Pos(self.0 - 1, self.1)
    }
    fn right(&self) -> Self {
        Pos(self.0 + 1, self.1)
    }
    fn manhattan(&self, other: &Pos) -> usize {
        distance(self.0, other.0) + distance(self.1, other.1)
    }
}

impl Ord for Pos {
    fn cmp(&self, other: &Pos) -> Ordering {
        match self.1.cmp(&other.1) {
            Ordering::Equal => self.0.cmp(&other.0),
            i => i,
        }
    }
}

impl PartialOrd for Pos {
    fn partial_cmp(&self, other: &Pos) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}

#[derive(Debug, Eq)]
struct Score {
    pos: Pos,
    g: usize,
}

impl Score {
    fn new(pos: Pos, g: usize) -> Self {
        Self { pos: pos, g: g }
    }
}

impl Ord for Score {
    fn cmp(&self, other: &Score) -> Ordering {
        match self.g.cmp(&other.g) {
            Ordering::Equal => self.pos.cmp(&other.pos).reverse(),
            i => i.reverse(),
        }
    }
}

impl PartialOrd for Score {
    fn partial_cmp(&self, other: &Score) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}

impl PartialEq for Score {
    fn eq(&self, other: &Score) -> bool {
        self.pos == other.pos
    }
}

#[derive(Clone)]
struct Entity {
    pos: Pos,
    elf: bool,
    health: isize,
    attack: isize,
}

impl Entity {
    fn new(pos: Pos, elf: bool) -> Self {
        Self {
            pos: pos,
            elf: elf,
            health: 200,
            attack: 3,
        }
    }

    fn dead(&self) -> bool {
        self.health < 1
    }

    fn damage(&mut self, attack: isize) -> bool {
        self.health -= attack;
        self.dead()
    }
}

struct Game<'a> {
    elfs: HashMap<Pos, Entity>,
    goblins: HashMap<Pos, Entity>,
    board: &'a Vec<Vec<char>>,
}

impl<'a> Game<'a> {
    fn new(
        elfs: HashMap<Pos, Entity>,
        goblins: HashMap<Pos, Entity>,
        board: &'a Vec<Vec<char>>,
    ) -> Self {
        Self {
            elfs: elfs,
            goblins: goblins,
            board: board,
        }
    }

    fn adjacent_enemies(&self, pos: Pos, elf: bool) -> Vec<(Pos, isize)> {
        let enemies = if elf { &self.goblins } else { &self.elfs };
        let mut adjacent = Vec::new();
        if let Some(i) = enemies.get(&pos.up()) {
            adjacent.push((pos.up(), i.health));
        }
        if let Some(i) = enemies.get(&pos.left()) {
            adjacent.push((pos.left(), i.health));
        }
        if let Some(i) = enemies.get(&pos.right()) {
            adjacent.push((pos.right(), i.health));
        }
        if let Some(i) = enemies.get(&pos.down()) {
            adjacent.push((pos.down(), i.health));
        }
        adjacent
    }

    fn is_free(&self, pos: Pos) -> bool {
        !self.elfs.get(&pos).is_some()
            && !self.goblins.get(&pos).is_some()
            && self.board[pos.1][pos.0] != '#'
    }

    fn adjacent_free(&self, pos: Pos) -> Vec<Pos> {
        let mut adjacent = Vec::new();
        if self.is_free(pos.up()) {
            adjacent.push(pos.up());
        }
        if self.is_free(pos.left()) {
            adjacent.push(pos.left());
        }
        if self.is_free(pos.right()) {
            adjacent.push(pos.right());
        }
        if self.is_free(pos.down()) {
            adjacent.push(pos.down());
        }
        adjacent
    }

    fn path(&self, start: Pos, goal: HashSet<Pos>) -> Pos {
        let mut closed = HashSet::new();
        let mut openset = HashSet::new();
        let mut open = BinaryHeap::new();
        for i in goal.into_iter() {
            open.push(Score::new(i, 0));
        }
        while let Some(k) = open.pop() {
            if start.manhattan(&k.pos) < 2 {
                return k.pos;
            }
            openset.remove(&k.pos);
            closed.insert(k.pos);
            for i in self.adjacent_free(k.pos) {
                if !closed.contains(&i) && !openset.contains(&i) {
                    open.push(Score::new(i, k.g + 1));
                    openset.insert(i);
                }
            }
        }
        start
    }

    fn move_entity(&self, pos: Pos, elf: bool) -> Pos {
        let enemies = if elf { &self.goblins } else { &self.elfs };
        let mut in_range = HashSet::new();
        for &p in enemies.keys() {
            for i in self.adjacent_free(p) {
                in_range.insert(i);
            }
        }
        self.path(pos, in_range)
    }

    fn tick_entity(&self, e: &Entity) -> (Pos, Option<Pos>) {
        let mut adjacent = self.adjacent_enemies(e.pos, e.elf);
        let mut next = e.pos;
        if adjacent.len() == 0 {
            next = self.move_entity(e.pos, e.elf);
            adjacent = self.adjacent_enemies(next, e.elf);
        }
        (
            next,
            match adjacent.iter().min_by_key(|i| i.1) {
                Some(i) => Some(i.0),
                None => None,
            },
        )
    }

    fn tick(&mut self) -> bool {
        let mut all = self
            .elfs
            .values()
            .chain(self.goblins.values())
            .map(|i| (i.elf, i.pos))
            .collect::<Vec<_>>();
        all.sort_by_key(|i| i.1);
        for (is_elf, i) in all.into_iter() {
            let entity = match if is_elf {
                self.elfs.get(&i)
            } else {
                self.goblins.get(&i)
            } {
                Some(i) => i,
                None => continue,
            };
            if self.elfs.len() == 0 || self.goblins.len() == 0 {
                return false;
            }

            let (next, target) = self.tick_entity(entity);
            let attack = if is_elf {
                let mut entity = self.elfs.remove(&i).expect("Failed to get elf");
                let a = entity.attack;
                entity.pos = next;
                self.elfs.insert(next, entity);
                a
            } else {
                let mut entity = self.goblins.remove(&i).expect("Failed to get goblin");
                let a = entity.attack;
                entity.pos = next;
                self.goblins.insert(next, entity);
                a
            };
            match target {
                Some(k) => {
                    let enemy = if is_elf {
                        self.goblins
                            .get_mut(&k)
                            .expect("Failed to get enemy goblin")
                    } else {
                        self.elfs.get_mut(&k).expect("Failed to get enemy elf")
                    };
                    if enemy.damage(attack) {
                        if is_elf {
                            self.goblins.remove(&k).expect("Failed to remove goblin");
                        } else {
                            self.elfs.remove(&k).expect("Failed to remove elf");
                        }
                    }
                }
                None => (),
            }
        }
        true
    }

    fn print(&self) {
        let mut b = self.board.clone();
        for k in self.elfs.keys() {
            b[k.1][k.0] = 'E';
        }
        for k in self.goblins.keys() {
            b[k.1][k.0] = 'G';
        }
        for i in b.iter() {
            println!("{}", i.iter().collect::<String>());
        }
    }
}

fn main() {
    let file = File::open(PUZZLEINPUT).expect("Failed to open file");
    let reader = BufReader::new(file);

    let (elfs, goblins, board) = {
        let mut elfs = HashMap::new();
        let mut goblins = HashMap::new();
        let mut board = Vec::new();
        for (y, line) in reader.lines().enumerate() {
            let mut k = line
                .expect("Failed to read line")
                .chars()
                .collect::<Vec<_>>();
            for (x, i) in k.iter_mut().enumerate() {
                match i {
                    'E' => {
                        elfs.insert(Pos(x, y), Entity::new(Pos(x, y), true));
                        *i = '.';
                    }
                    'G' => {
                        goblins.insert(Pos(x, y), Entity::new(Pos(x, y), false));
                        *i = '.';
                    }
                    _ => (),
                }
            }
            board.push(k);
        }
        (elfs, goblins, board)
    };

    let mut game = Game::new(elfs.clone(), goblins.clone(), &board);
    let mut round = 0;
    while game.tick() {
        round += 1;
    }
    game.print();
    let total_health = game.elfs.values().map(|i| i.health).sum::<isize>()
        + game.goblins.values().map(|i| i.health).sum::<isize>();
    println!("{}", round * total_health);
}
