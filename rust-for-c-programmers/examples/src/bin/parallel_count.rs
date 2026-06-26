// Project 2: a concurrent word-frequency counter.
//
// We split a body of text into chunks, hand each chunk to a worker thread,
// and merge the per-worker results into one shared table. The number of
// workers is bounded. Every thread is joined, so the program always ends.
//
// Two safe sharing tools are shown:
//   * std::sync::mpsc  -- workers SEND their partial maps to the main thread.
//   * std::thread::scope -- lets threads borrow stack data without 'static.
//
// The merge order does not change the totals, so the output is deterministic.

use std::collections::HashMap;
use std::sync::mpsc;
use std::thread;

/// Sample input. In a real tool this would come from a file or stdin.
const TEXT: &str = "\
the quick brown fox
the lazy dog
the fox jumps over the lazy dog
brown fox brown dog
the the the quick brown fox";

/// How many worker threads to run. Bounded on purpose.
const WORKERS: usize = 4;

/// Count words in one slice of lines into a fresh map.
fn count_chunk(lines: &[&str]) -> HashMap<String, u64> {
    let mut local = HashMap::new();
    for line in lines {
        for word in line.split_whitespace() {
            *local.entry(word.to_string()).or_insert(0) += 1;
        }
    }
    local
}

/// Split the lines into at most `n` roughly equal chunks.
fn split_into_chunks<'a>(lines: &'a [&'a str], n: usize) -> Vec<&'a [&'a str]> {
    if lines.is_empty() {
        return Vec::new();
    }
    let n = n.min(lines.len()).max(1);
    let chunk_size = lines.len().div_ceil(n);
    lines.chunks(chunk_size).collect()
}

/// Run the parallel count and return the merged table.
fn parallel_word_count(text: &str) -> HashMap<String, u64> {
    let lines: Vec<&str> = text.lines().collect();
    let chunks = split_into_chunks(&lines, WORKERS);

    let mut totals: HashMap<String, u64> = HashMap::new();

    // A channel: workers send partial maps, the main thread receives them.
    let (tx, rx) = mpsc::channel::<HashMap<String, u64>>();

    // `scope` guarantees all spawned threads finish before it returns, so the
    // threads may safely borrow `chunks` (no 'static requirement, no Arc).
    thread::scope(|s| {
        for chunk in &chunks {
            let tx = tx.clone();
            let chunk = *chunk; // copy the &[&str] slice handle, not the data
            s.spawn(move || {
                let local = count_chunk(chunk);
                // If the receiver is gone we just drop the result; ignore error.
                let _ = tx.send(local);
            });
        }
        // Drop the original sender so the receiver loop can end once all
        // worker clones are dropped at the end of their threads.
        drop(tx);

        // Merge results as they arrive. The loop ends when every sender is gone.
        for partial in rx {
            for (word, n) in partial {
                *totals.entry(word).or_insert(0) += n;
            }
        }
    });

    totals
}

fn main() {
    let totals = parallel_word_count(TEXT);

    // Sort for stable, human-friendly output: by count desc, then by word.
    let mut rows: Vec<(&String, &u64)> = totals.iter().collect();
    rows.sort_by(|a, b| b.1.cmp(a.1).then(a.0.cmp(b.0)));

    println!("word frequencies ({} unique words):", rows.len());
    for (word, count) in rows {
        println!("{count:>3}  {word}");
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn counts_are_correct_and_deterministic() {
        let totals = parallel_word_count(TEXT);
        assert_eq!(totals["the"], 7);
        assert_eq!(totals["fox"], 4);
        assert_eq!(totals["brown"], 4);
        assert_eq!(totals["dog"], 3);
        assert_eq!(totals["quick"], 2);
        assert_eq!(totals["jumps"], 1);
    }

    #[test]
    fn empty_input_yields_empty_table() {
        let totals = parallel_word_count("");
        assert!(totals.is_empty());
    }

    #[test]
    fn merging_order_does_not_change_totals() {
        // Running twice must give identical results (no data race, no nondeterminism).
        let a = parallel_word_count(TEXT);
        let b = parallel_word_count(TEXT);
        assert_eq!(a, b);
    }
}
