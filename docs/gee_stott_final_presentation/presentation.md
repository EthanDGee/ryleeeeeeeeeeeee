---
marp: true
theme: default
paginate: true
style: |
  section {
    background-color: #f5f5f5;
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    font-size: 23px;
  }
  section.lead {
    background: linear-gradient(135deg, #ccccff 0%, #a754fbff 100%);
    color: white;
    text-align: center;
  }
  section.lead h1 {
    font-size: 3em;
    text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
  }
  h1 {
    color: #2c3e50;
    border-bottom: 4px solid #667eea;
    padding-bottom: 10px;
  }
  h2 {
    color: #34495e;
    margin-top: 1em;
  }
  strong {
    color: #7777ea;
  }
  code {
    background-color: #2c3e50;
    color: #ecf0f1;
    padding: 2px 6px;
    border-radius: 3px;
  }
  pre {
    background-color: #2c3e50;
    border-radius: 8px;
    padding: 20px;
  }
  pre code {
    background-color: transparent;
  }
  ul, ol {
    line-height: 1.8;
  }
  li {
    margin: 0.5em 0;
  }
  blockquote {
    border-left: 5px solid #667eea;
    padding-left: 1em;
    font-style: italic;
    color: #555;
  }
  footer {
    color: #7f8c8d;
  }
  section::after {
    color: #7f8c8d;
    font-weight: bold;
  }
  .columns {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 1rem;
  }
---

<!-- _class: lead -->

# A Human-Like Chess Engine
## Rylee

**Ethan Gee & Nate Stott**

---

# Agenda

1. **Introduction & Motivation**
2. **Methodology**
3. **Experiments**
4. **Conclusions & Future Work**

---

# Introduction - Problem Definition

<div class="columns">
<div>

**The Maia Problem**
- Traditional engines play **optimally**, not **humanly**
- Predicting human moves != finding the best move
- **Goal**: Predict what a human player would actually choose

</div>
<div>

**The Rylee Problem**
- Maia requires **GPU & large model**
- Not edge-deployable
- Missing opening play
- **Goal**: Human-like play on **edge hardware**

</div>
</div>

---

# Introduction - Motivation

<div class="columns">
<div>

**Why human-like AI matters**
- Traditional engines are **too strong** to learn from
- Handicapping creates **unnatural play**
- Human-like engines = **realistic training partners**
- Example: Chess students can practice with Rylee to **experience human-style mistakes and strategies**, improving learning
- Broader applications: Education, coaching, collaborative decision-making

</div>
<div>

**Why Rylee extends Maia**
- **Edge deployment** (Raspberry Pi, Jetson)
- **Includes openings** (first 10 moves)
- **Unified model** (700-2500 ELO)
- **Human-like behavior at all levels** allows consistent training experience

</div>
</div>

---

# Methodology - Proposed Solution

<img src="./figures/high_level.png" width="1000">

Train on millions of human chess positions to predict human moves using a lightweight CNN.

---

# Methodology - Theories

<div class="columns">
<div>

**Why CNNs work well for chess**
- Chess patterns are **spatially local** (pawn structures, king safety)
- CNNs excel at **spatial pattern recognition**

**Model size vs performance**
- Smaller models enable **edge deployment**
- Performance doesn't scale linearly with size

</div>
<div>

**Rylee = Maia**
- Both predict **human moves** using supervised learning

**Rylee != Maia**
- **Smaller model** (edge-deployable vs GPU-only)
- **Includes openings** (first 10 moves)
- **Unified model** (700-2500 ELO vs 9 separate models)

</div>
</div>

---

# Methodology - Data Pipeline

<div class="columns">
<div>

<img src="./figures/data_pipeline.png" width="600">

</div>
<div>

- **Download** .zst files from Lichess
- **Extract** PGNs
- **Split** into individual games
- **Convert** to board state snapshots
- **Extract** ELO & result metadata
- **Encode** as one-hot tensors

</div>
</div>

---

# Methodology - Neural Network

<div class="columns">
<div>

<img src="./figures/nn_architecture.png" width="600">

</div>
<div>

- **Input**: Board(12x8x8) + Metadata(4)
- **Conv Layers**: 6x64 filters @ 8x8, ReLU
- **Fully Connected**: 4100 -> 512 -> 32
- **Output Heads**: Move (2104) + Auxiliary (2104)
- **Loss**: CrossEntropy (moves) + BCE (valid moves)
- **Optimizer**: Adam
- **Hyperparameter Search**: Random search

</div>
</div>

---

# Experiments - Dataset

<div class="columns">
<div>

- **Source:** Lichess Open Database
- **Games:** 15,000 human-rated games
- **Snapshots:** 1 million board positions
  - including openings
  - including all game types
- **Action Space:** 2,104 legal move classes
- **Time Span:** January 2013

</div>
<div>

| Split      | Percentage | Snapshots      |
|------------|------------|----------------|
| Training   | 80%        | **800,000**    |
| Validation | 10%        | **100,000**    |
| Test       | 10%        | **100,000**    |

<img src="./figures/elo-distribution.png" width="600">

</div>
</div>

---

# Experiments - Baselines

| Method            | Description                      |
|-------------------|----------------------------------|
| **Random**        | Random legal move selection      |
| **Random Forest** | Classic ML with handcrafted features |
| **Stockfish 15**  | Traditional chess engine         |
| **Leela 4200**    | Neural chess engine              |
| **Maia1 1500**    | Human-aligned prediction model   |

---

# Experiments - Evaluation Metrics

- **Top-1 Accuracy**: Predicted move matches actual human move
- **Top-5 Accuracy**: Actual move in top 5 predictions

---

# Experiments - Comparisons

<div class="columns">
<div>

| Method            | Top-1 Accuracy |
|-------------------|----------------|
| **Random**        | **6%**         |
| **Random Forest** | **13%**        |
| **Stockfish 15**  | **40%**        |
| **Leela 4200**    | **44%**        |
| **Maia1 1500**    | **51%**        |
| **Rylee (Ours)**  | **25%**        |

**Rylee - Key Differences**

- Rylee = 800,000 parameters
- Trained on Raspberry pi and they trained on

</div>
<div>

- Maia = 25 million parameters
- Trained on 2 A100 80Gb GPUs
- No filtering by game type (e.g., blitz/classical) - aims to capture broader human play patterns.
- No Elo filtering - unlike Maia, we include games with mixed skill levels to better reflect general human behavior.
- No data augmentation currently used.
- Dataset is **far smaller**: ~15,000 games vs. Maia’s 169 million games 1.9 billion snapshots.
- Model size is **~1/5** that of Maia.

</div>
</div>

---

# Conclusions - Discussions

<div class="columns">
<div>

| Metric          | Training | Validation |
|-----------------|----------|------------|
| **Loss**        | 0.0152   | 0.0164     |
| **Top-1 Accuracy** | 28%    | 25%      |
| **Top-5 Accuracy** | 53%    | 51%      |

- Strong generalization between training and validation metrics
- Model captures key patterns of human decision-making
- Rylee required less than 1 day of preprocessing and 2-3 days of training

</div>
<div>

- Maia required 8 days of preprocessing and 3-4 weeks of training

<img src="./figures/epoch_curves.png" width="400">

</div>
</div>

---

# Conclusions - Future Work

**Model Improvements**
- Add data augmentation (flips, rotations) to improve robustness.
- Fine-tune rating-specific models for better skill-level alignment.

**Additional Features**
- **ELO Prediction:** Estimate player rating from move patterns.
- **Human vs Bot Discriminator:** Detect engine-like play.
- **Blunder Detection:** Identify major mistakes for analysis.

---

# References

<div class="columns">
<div>

**Primary Works**
- McIlroy-Young et al. (2020). "Aligning Superhuman AI with Human Behavior: Chess as a Model System." KDD 2020.
- Tang et al. (2024). "Maia-2: A Unified Model for Human-AI Alignment in Chess." NeurIPS 2024.
- McIlroy-Young et al. (2021). "Detecting Individual Decision-Making Style: Exploring Behavioral Stylometry in Chess." NeurIPS 2021.

</div>
<div>

**Data & Tools**
- Lichess Open Database: https://database.lichess.org/
- Stockfish Chess Engine: https://stockfishchess.org/
- Leela Chess Zero: https://lczero.org/
- Maia Chess Project: https://maiachess.com/
- PyTorch (Paszke et al., 2019): https://pytorch.org/
- python-chess library (Moskopp, 2014): https://github.com/niklasf/python-chess

</div>
</div>

---

<!-- _class: lead -->

# Questions?

**Rylee**: Human-Like Chess Engine

Ethan Gee & Nate Stott
