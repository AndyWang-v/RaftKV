//! Project 3: a small generic data structure with tests.
//!
//! This crate gives a generic `Stack<T>` and, on top of it, a tiny
//! reverse-polish-notation (RPN) calculator. It exercises generics, traits,
//! `Option`, and `Result`, and it ships with both `#[test]` unit tests and a
//! documentation test (the example on `Stack` below is compiled and run).
//!
//! ```
//! use rustbook_examples::Stack;
//!
//! let mut s: Stack<i32> = Stack::new();
//! s.push(10);
//! s.push(20);
//! assert_eq!(s.len(), 2);
//! assert_eq!(s.peek(), Some(&20));
//! assert_eq!(s.pop(), Some(20));
//! assert_eq!(s.pop(), Some(10));
//! assert_eq!(s.pop(), None);
//! assert!(s.is_empty());
//! ```

use std::fmt;

/// A growable last-in, first-out stack of `T`, backed by a `Vec<T>`.
///
/// This is generic over any element type, the same way a C version would use
/// `void *` plus a size — but here the type is checked at compile time and the
/// code is monomorphized, so there is no runtime cost.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Stack<T> {
    items: Vec<T>,
}

impl<T> Stack<T> {
    /// Create an empty stack.
    pub fn new() -> Self {
        Stack { items: Vec::new() }
    }

    /// Create an empty stack with room for `cap` items before reallocating.
    pub fn with_capacity(cap: usize) -> Self {
        Stack {
            items: Vec::with_capacity(cap),
        }
    }

    /// Push a value onto the top.
    pub fn push(&mut self, value: T) {
        self.items.push(value);
    }

    /// Remove and return the top value, or `None` if the stack is empty.
    pub fn pop(&mut self) -> Option<T> {
        self.items.pop()
    }

    /// Look at the top value without removing it.
    pub fn peek(&self) -> Option<&T> {
        self.items.last()
    }

    /// Number of items on the stack.
    pub fn len(&self) -> usize {
        self.items.len()
    }

    /// True when the stack has no items.
    pub fn is_empty(&self) -> bool {
        self.items.is_empty()
    }
}

// A common Rust idiom: implement `Default` so `Stack::default()` works and so
// the type fits APIs that need a default value.
impl<T> Default for Stack<T> {
    fn default() -> Self {
        Stack::new()
    }
}

// Let any iterator of `T` build a stack with `.collect()`.
impl<T> FromIterator<T> for Stack<T> {
    fn from_iter<I: IntoIterator<Item = T>>(iter: I) -> Self {
        Stack {
            items: iter.into_iter().collect(),
        }
    }
}

/// The reason an RPN expression could not be evaluated.
#[derive(Debug, Clone, PartialEq, Eq)]
pub enum RpnError {
    /// A token was neither a number nor a known operator.
    BadToken(String),
    /// An operator ran but the stack did not have two operands.
    NotEnoughOperands,
    /// Division (or remainder) by zero.
    DivideByZero,
    /// The expression finished with more or fewer than one value left.
    LeftoverOperands(usize),
}

// Implementing `Display` makes `RpnError` print nicely and (with the blanket
// impl in std) usable as a boxed error. This is the trait-driven error style.
impl fmt::Display for RpnError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            RpnError::BadToken(t) => write!(f, "bad token: {t}"),
            RpnError::NotEnoughOperands => write!(f, "not enough operands"),
            RpnError::DivideByZero => write!(f, "divide by zero"),
            RpnError::LeftoverOperands(n) => {
                write!(f, "expression left {n} values on the stack")
            }
        }
    }
}

impl std::error::Error for RpnError {}

/// Evaluate a reverse-polish-notation expression of `i64` integers.
///
/// Supported operators: `+`, `-`, `*`, `/`, `%`. Tokens are separated by
/// whitespace. Returns the single result or an [`RpnError`].
///
/// ```
/// use rustbook_examples::eval_rpn;
///
/// assert_eq!(eval_rpn("2 3 +"), Ok(5));
/// assert_eq!(eval_rpn("4 2 - 3 *"), Ok(6)); // (4 - 2) * 3
/// assert!(eval_rpn("1 0 /").is_err());
/// ```
pub fn eval_rpn(expr: &str) -> Result<i64, RpnError> {
    let mut stack: Stack<i64> = Stack::new();

    for token in expr.split_whitespace() {
        if let Ok(n) = token.parse::<i64>() {
            stack.push(n);
            continue;
        }

        // An operator: pop the two operands (note the order: b is on top).
        let b = stack.pop().ok_or(RpnError::NotEnoughOperands)?;
        let a = stack.pop().ok_or(RpnError::NotEnoughOperands)?;

        let result = match token {
            "+" => a + b,
            "-" => a - b,
            "*" => a * b,
            "/" => {
                if b == 0 {
                    return Err(RpnError::DivideByZero);
                }
                a / b
            }
            "%" => {
                if b == 0 {
                    return Err(RpnError::DivideByZero);
                }
                a % b
            }
            other => return Err(RpnError::BadToken(other.to_string())),
        };
        stack.push(result);
    }

    match stack.len() {
        1 => Ok(stack.pop().unwrap()),
        n => Err(RpnError::LeftoverOperands(n)),
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn push_pop_peek() {
        let mut s = Stack::new();
        assert!(s.is_empty());
        s.push("a");
        s.push("b");
        assert_eq!(s.peek(), Some(&"b"));
        assert_eq!(s.len(), 2);
        assert_eq!(s.pop(), Some("b"));
        assert_eq!(s.pop(), Some("a"));
        assert_eq!(s.pop(), None);
    }

    #[test]
    fn works_for_different_types() {
        // The same generic code handles i32 and String with no extra work.
        let mut ints: Stack<i32> = Stack::new();
        ints.push(1);
        ints.push(2);
        assert_eq!(ints.pop(), Some(2));

        let mut words: Stack<String> = Stack::new();
        words.push("hello".to_string());
        assert_eq!(words.peek().map(|s| s.as_str()), Some("hello"));
    }

    #[test]
    fn collect_into_stack() {
        let s: Stack<i32> = (1..=3).collect();
        assert_eq!(s.len(), 3);
        assert_eq!(s.peek(), Some(&3));
    }

    #[test]
    fn rpn_basic() {
        assert_eq!(eval_rpn("2 3 +"), Ok(5));
        assert_eq!(eval_rpn("10 2 /"), Ok(5));
        assert_eq!(eval_rpn("7 3 %"), Ok(1));
        assert_eq!(eval_rpn("2 3 4 * +"), Ok(14));
    }

    #[test]
    fn rpn_errors() {
        assert_eq!(eval_rpn("1 +"), Err(RpnError::NotEnoughOperands));
        assert_eq!(eval_rpn("1 2"), Err(RpnError::LeftoverOperands(2)));
        assert_eq!(eval_rpn("1 0 /"), Err(RpnError::DivideByZero));
        assert_eq!(
            eval_rpn("1 2 foo"),
            Err(RpnError::BadToken("foo".to_string()))
        );
    }

    #[test]
    fn rpn_error_displays() {
        let msg = eval_rpn("1 0 /").unwrap_err().to_string();
        assert_eq!(msg, "divide by zero");
    }
}
