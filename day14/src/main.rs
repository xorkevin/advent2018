const PUZZLEINPUT: usize = 540391;
const PUZZLEINPUT2: &[u8] = &[5, 4, 0, 3, 9, 1];

fn split(b: u8) -> (u8, u8) {
    (b / 10, b % 10)
}

fn compare(state: &[u8], target: &[u8]) -> usize {
    let sl = state.len();
    let tl = target.len();
    if sl < tl + 1 {
        return 0;
    }
    let l = sl - tl;
    if &state[l - 1..sl - 1] == target || &state[l..] == target {
        return l;
    }
    return 0;
}

fn main() {
    part1();
    part2();
}

fn part1() {
    let mut count1 = 0;
    let mut count2 = 1;
    let mut state = vec![3, 7];
    while state.len() < PUZZLEINPUT + 10 {
        let (a, b) = split(state[count1] + state[count2]);
        if a > 0 {
            state.push(a);
        }
        state.push(b);
        count1 = (count1 + state[count1] as usize + 1) % state.len();
        count2 = (count2 + state[count2] as usize + 1) % state.len();
    }
    println!(
        "{}",
        state[PUZZLEINPUT..PUZZLEINPUT + 10]
            .iter()
            .map(|&x| (x + '0' as u8) as char)
            .collect::<String>()
    );
}

fn part2() {
    let mut count1 = 0;
    let mut count2 = 1;
    let mut state = vec![3, 7];

    let mut k = 0;
    while k == 0 {
        let (a, b) = split(state[count1] + state[count2]);
        if a > 0 {
            state.push(a);
        }
        state.push(b);
        count1 = (count1 + state[count1] as usize + 1) % state.len();
        count2 = (count2 + state[count2] as usize + 1) % state.len();
        k = compare(&state, PUZZLEINPUT2);
    }
    println!("{}", k);
}
