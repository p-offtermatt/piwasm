use std::collections::{HashMap, HashSet};

pub fn require(cond: bool) -> bool {
    cond
}

#[cfg(test)]
mod requireTest {
    use super::*;

    #[test]
    fn test() {
        assert!(require(4 > 3));
        assert!(!require(false));
    }
}

pub fn requires(cond: bool, error: &str) -> &str {
    if cond {
        ""
    } else {
        error
    }
}

#[cfg(test)]
mod requiresTest {
    use super::*;

    #[test]
    fn test() {
        assert!(requires(4 > 3, "4 > 3") == "");
        assert!(requires(4 < 3, "false: 4 < 3") == "false: 4 < 3");
    }
}

pub fn max(i: i32, j: i32) -> i32 {
    if i > j {
        i
    } else {
        j
    }
}

#[cfg(test)]
mod maxTest {
    use super::*;

    #[test]
    fn test() {
        assert!(max(3, 4) == 4);
        assert!(max(6, 3) == 6);
        assert!(max(10, 10) == 10);
        assert!(max(-3, -5) == -3);
        assert!(max(-5, -3) == -3);
    }
}

pub fn abs(i: i32) -> i32 {
    i.abs()
}

#[cfg(test)]
mod test_abs {
    use super::*;

    #[test]
    fn test() {
        assert!(abs(3) == 3);
        assert!(abs(-3) == 3);
        assert!(abs(0) == 0);
    }
}

// FIXME(romain): we probably need to special case this function,
//                as we can't generically infer the bounds nor
//                can we translate it directly from Quint
pub fn setRemove<T: std::cmp::Eq + std::hash::Hash + std::clone::Clone>(
    set: &HashSet<T>,
    elem: &T,
) -> HashSet<T> {
    let mut new_set = set.clone();
    new_set.remove(elem);
    new_set
}

#[cfg(test)]
mod setRemoveTest {
    use super::*;

    #[test]
    fn test() {
        let mut a = std::collections::HashSet::new();
        a.insert(2);
        a.insert(3);
        a.insert(4);
        let mut b = std::collections::HashSet::new();
        b.insert(2);
        b.insert(4);
        assert!(b == setRemove(&a, &3));
        let mut c = std::collections::HashSet::new();
        assert!(c == setRemove(&c, &3));
    }
}

// FIXME(romain): we probably also need to special case this function
pub fn has<K: std::cmp::Eq + std::hash::Hash, V>(__map: &HashMap<K, V>, __key: &K) -> bool {
    __map.contains_key(__key)
}

#[cfg(test)]
mod hasTest {
    use super::*;

    #[test]
    fn test() {
        let mut a = std::collections::HashMap::new();
        a.insert(2, 3);
        a.insert(4, 5);
        assert!(has(&a, &2));
        assert!(!has(&a, &6));
    }
}

// FIXME(romain): we probably also need to special case this function
pub fn getOrElse<K: std::cmp::Eq + std::hash::Hash + std::clone::Clone, V: std::clone::Clone>(
    __map: &HashMap<K, V>,
    __key: &K,
    __default: V,
) -> V {
    if __map.contains_key(__key) {
        __map.get(__key).unwrap().clone()
    } else {
        __default
    }
}

#[cfg(test)]
mod getOrElseTest {
    use super::*;

    #[test]
    fn test() {
        let mut a = std::collections::HashMap::new();
        a.insert(2, 3);
        a.insert(4, 5);
        assert!(getOrElse(&a, &2, 0) == 3);
        assert!(getOrElse(&a, &7, 11) == 11);
    }
}

// FIXME(romain): we probably also need to special case this function
pub fn mapRemove<K: std::cmp::Eq + std::hash::Hash + std::clone::Clone, V: std::clone::Clone>(
    __map: &HashMap<K, V>,
    __key: &K,
) -> HashMap<K, V> {
    let mut new_map = __map.clone();
    new_map.remove(__key);
    new_map
}

#[cfg(test)]
mod mapRemoveTest {
    use std::collections::HashMap;

    use super::*;

    #[test]
    fn test() {
        let mut a = HashMap::new();
        a.insert(3, 4);
        a.insert(5, 6);
        a.insert(7, 8);
        let mut b = HashMap::new();
        b.insert(3, 4);
        b.insert(7, 8);
        assert!(b == mapRemove(&a, &5));
        // let mut c = HashMap::new();
        // assert!(c == mapRemove(&c, &3));
    }
}

// FIXME(romain): we probably also need to special case this function
pub fn mapRemoveAll<K: std::cmp::Eq + std::hash::Hash + std::clone::Clone, V: std::clone::Clone>(
    __map: &HashMap<K, V>,
    __keys: &HashSet<K>,
) -> HashMap<K, V> {
    let mut new_map = __map.clone();
    for key in __keys {
        new_map.remove(key);
    }
    new_map
}

#[cfg(test)]
mod mapRemoveAllTest {
    use std::collections::{HashMap, HashSet};

    use super::*;

    #[test]
    fn test() {
        let mut a = HashMap::new();
        a.insert(3, 4);
        a.insert(5, 6);
        a.insert(7, 8);
        let mut keys = HashSet::new();
        keys.insert(5);
        keys.insert(7);
        let mut b = HashMap::new();
        b.insert(3, 4);
        assert!(b == mapRemoveAll(&a, &keys));
        let mut keys = HashSet::new();
        keys.insert(5);
        keys.insert(99999);
        let mut c = HashMap::new();
        c.insert(3, 4);
        c.insert(7, 8);
        assert!(c == mapRemoveAll(&a, &keys));
    }
}
