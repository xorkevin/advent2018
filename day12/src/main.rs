use std::collections::HashMap;

const PUZZLEINPUT: &str = "#.#..#..###.###.#..###.#####...########.#...#####...##.#....#.####.#.#..#..#.#..###...#..#.#....##.";
const ITERATIONS: i64 = 50000000000;

fn exec_rule(ind: usize, state: &Vec<char>, rules: &HashMap<String, char>) -> char {
    let left = ind - 2;
    let right = ind + 3;
    let s = state[left..right].iter().collect::<String>();
    if let Some(&val) = rules.get(&s) {
        val
    } else {
        state[ind]
    }
}

fn score(zero_ind: i64, state: &Vec<char>) -> i64 {
    state.iter().enumerate().fold(0, |acc, (n, &i)| {
        if i == '#' {
            acc + n as i64 - zero_ind
        } else {
            acc
        }
    })
}

fn next_state(
    state: Vec<char>,
    mut next: Vec<char>,
    rules: &HashMap<String, char>,
) -> (Vec<char>, Vec<char>) {
    next[0] = '.';
    next[1] = '.';
    next[state.len() - 2] = '.';
    next[state.len() - 1] = '.';
    for i in 2..state.len() - 2 {
        next[i] = exec_rule(i, &state, rules);
    }
    (next, state)
}

fn main() {
    let rules = {
        let k = vec![
            ("#.###", '.'),
            ("###.#", '#'),
            (".##..", '.'),
            ("..###", '.'),
            ("..##.", '.'),
            ("##...", '#'),
            ("###..", '#'),
            (".#...", '#'),
            ("##..#", '#'),
            ("#....", '.'),
            (".#.#.", '.'),
            ("####.", '.'),
            ("#.#..", '.'),
            ("#.#.#", '.'),
            ("#..##", '#'),
            (".####", '#'),
            ("...##", '.'),
            ("#..#.", '#'),
            (".#.##", '#'),
            ("..#.#", '#'),
            ("##.#.", '#'),
            ("#.##.", '#'),
            ("#####", '.'),
            ("..#..", '#'),
            ("....#", '.'),
            ("##.##", '.'),
            (".###.", '#'),
            (".....", '.'),
            ("...#.", '#'),
            (".##.#", '.'),
            ("#...#", '.'),
            (".#..#", '#'),
        ];
        let mut rules = HashMap::new();
        for (k, v) in k.into_iter() {
            rules.insert(k.to_string(), v);
        }
        rules
    };
    let (mut state, zero_ind) = {
        let zero_ind = PUZZLEINPUT.len();
        let mut state = vec!['.'; PUZZLEINPUT.len() * 16];
        state[zero_ind..zero_ind + PUZZLEINPUT.len()]
            .copy_from_slice(PUZZLEINPUT.chars().collect::<Vec<_>>().as_slice());
        (state, zero_ind as i64)
    };

    let mut prev_delta = 0;
    let mut prev_score = score(zero_ind, &state);
    let mut prev_state = vec!['.'; state.len()];
    for i in 0..256 {
        let (s1, p1) = next_state(state, prev_state, &rules);
        state = s1;
        prev_state = p1;
        let score = score(zero_ind, &state);
        let delta = score - prev_score;
        if i == 19 {
            println!("{}", score);
        }
        if delta == prev_delta {
            println!("{}", score + (ITERATIONS - i - 1) * delta);
            return;
        }
        prev_score = score;
        prev_delta = delta;
    }
}
