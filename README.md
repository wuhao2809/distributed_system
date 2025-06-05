# MIT 6.5840 (Distributed Systems) — Lab Solutions

This repository contains my personal solutions to the [MIT 6.5840: Distributed Systems](https://pdos.csail.mit.edu/6.5840/) labs.

## 📁 Structure

- `Makefile`: Top-level make targets for each lab.
- `src/`: Contains all lab source code and test harnesses.
  - `raft/`: Implementation of the Raft consensus algorithm.
  - `kvraft/`: Key-Value store built on top of Raft.
  - `mr/`: MapReduce implementation.

## 🧪 Labs Overview

| Lab   | Description                        |
| ----- | ---------------------------------- |
| Lab 1 | MapReduce                          |
| Lab 2 | Raft Consensus Algorithm (Part 1)  |
| Lab 3 | Raft Log Compaction & Snapshotting |
| Lab 4 | Fault-Tolerant Key/Value Server    |

## ⚙️ Building & Testing

To compile and run tests:

```bash
cd src
make lab1       # Replace lab1 with lab2, lab3, or lab4 as needed
```

## Status

- [ ] Lab 1: In progress
- [ ] Lab 2: Not started
- [ ] Lab 3: Not started
- [ ] Lab 4: Not started

## Disclaimer

This repository is intended for **personal learning and documentation only**. Please do not copy the code for academic submissions. Follow the [MIT Honor Code](https://integrity.mit.edu/) and use responsibly.
