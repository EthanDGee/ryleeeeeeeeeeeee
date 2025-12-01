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

# Creating a Deployable Human-Like Chess Engine to Enhance the Learning Experience

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

**Maia Problems**
- Traditional chess engines maximize the chances of winning
- Predicting a move a human would play does not mean finding the best move
- **Goal**: Replicate human play

</div>
<div>

**Rylee Problems**
- Maia requires **large GPUs** to run and train
- Maia can only run on high end machines
- Maia was not trained to play chess openings
- **Goal**: Make a **practical & deployable#** version of Maia

</div>
</div>

---

# Introduction - Motivation

<div class="columns">
<div>

**Why human aligned AI matters**
- Traditional engines play chess differently making it difficult for humans to learn from
- Attenuating does not **mimic human play**
- Human aligned engines creates more realistic **training partners**
- Example: Chess students can practice with Rylee on their school chromebook to advance their chess skills
- Broader applications: Collaborative decision-making, Education, etc

</div>
<div>

**How Rylee extends Maia**
- **Edge deployment** Raspberry Pi, Chromebooks
- **Includes openings** first 10 moves
- **No game filtering** include all game types (classical, blitz, etc)
- **Unified model with a higher range** 700-2500 ELO

</div>
</div>

---

# Methodology - Proposed Solution

We believes we can still maintain similar performance, and add feautures to the MAIA model while signficantly reducing model size.

<img src="./figures/high_level.png" width="1000">

1. Pull data from Lichess
2. Preprocess games
3. Feed data into NN
4. Predict moves

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
- **Convert** to board snapshots
- **Extract** ELO and result metadata
- **Encode** board as 8x8x12 tensors

</div>
</div>

---

# Methodology - Theories


**Why CNNs work well for chess**
- Chess boards are **spatially related** (knight is better if its in the middle)
- Humans evaluate through pattern recognition
- CNNs excel at **spatial pattern recognition**

**Model size vs performance**
- Increasing model size has exponentially diminishing returns
- We are hoping we can decrease the size of the Maia model while still keeping high accuracy

---

# Methodology - Neural Network

<div class="columns">
<div>

<img src="./figures/nn_architecture.png" width="600">

</div>
<div>

- **Input**: Board(8x8x12) + Metadata(4)
- **Conv Layers**: 6 64x8x8 filters, ReLU
- **Fully Connected**: 4100 -> 512 -> 32
- **Output Heads**: Move (2104) + Auxiliary (2104)
  - 2104 is the number of legal moves
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
- **Games:** 15,167 human-rated games
- **Snapshots:** 1 million board states
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

| Baseline Model            | Description                      |
|-------------------|----------------------------------|
| **Random**        | Random legal move selection      |
| **Random Forest** | Nothing that simple should work that well - Ethan Gee |
| **Stockfish 15**  | Traditional chess engine         |
| **Leela 4200**    | Neural chess engine              |
| **Maia-1 1500**    | Human aligned prediction model   |

---

# Experiments - Architecture

### Small Fully Connected Model
A Small model that had a similar architecture to StockFish
- 8 fully connected layers of 32 neurons

### Convolutional Model 
- Combination of Convolution and fully connected to mirror human cognition

### Convolution with Auxillary Head
- Added an auxillary head that determines legal moves to instill better game understanding

---

# Experiments - Evaluation Metrics

- **Top-1 Accuracy**: Predicted move matches actual human move
- **Top-5 Accuracy**: Actual move in top 5 predictions. This is a good for a more generalizaed alignment.

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
| **Rylee**         | **25%**        |

- Rylee has 800,000 parameters
  - Trained on one Raspberry pi
- Maia has 25 million parameters
  - Trained on two A100 80Gb GPUs

</div>

<div>

- Rylee has 800,000 parameters vs MAIAs 25 Million
- No filtering by game type (classical, blitz, etc) to capture broader human play patterns
- We include games with mixed skill levels to better reflect general human behavior
- 15,000 games vs. Maia’s 169 million games
- MAIA was Trained on two A100 80Gb GPUs vs Rylee being trained on a Edge Device
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
- Model captures key human decision-making patterns
- Rylee required around 1.5 hours of preprocessing and 2-3 days of training

</div>
<div>

- Maia required 8 days of preprocessing and 3-4 weeks of training

<img src="./figures/epoch_curves.png" width="400">

</div>
</div>

---

# Conclusions - Future Work

**Model Improvements**
- Add data augmentation (board flips and rotations) to improve robustness
- Time parameter to better address time based decision making

**Additional Features**
- **ELO Prediction:** Estimate player rating from move patterns to quickly adapt to player skill
- **Human vs Bot Discriminator:** Detect engine-like play
- **Blunder Detection:** Identify major mistakes for analysis

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
