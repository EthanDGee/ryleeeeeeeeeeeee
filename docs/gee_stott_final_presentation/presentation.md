---
marp: true
theme: default
paginate: true
style: |
  section {
    background-color: #f5f5f5;
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
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

<!-- TODO: OVERALL PRESENTATION REQUIREMENTS:
     - Time limit: 10-12 minutes total
     - Practice timing to ensure you cover all sections properly
     - Add references slide at the end (before Questions slide) citing:
       * Maia paper
       * Lichess database
       * Any other sources used
     - Verify all font sizes are readable (current styling looks good)
     - Ensure all figures/charts are visible when presenting
     - Check that bullet points are concise and well-summarized
-->

<!-- _class: lead -->

# Rylee
## A Human-Like Chess Engine

**Ethan Gee & Nate Stott**

---

# Agenda

1. **Introduction & Motivation**
2. **Methodology**
3. **Experiments**
4. **Conclusions & Future Work**

---

# Introduction - Problem Definition

**Maia**
- Traditional chess engines (Stockfish, LC0) play **optimally**
- They don't play like humans
- Goal: Predict what move a **human** would make

**Rylee** - All the above AND
- Maia is a **large complex model** that takes a **lots of compute power** to train and run
- Goal: Get about the **same accuracy** but train and run on a **raspberry pi**

---

# Introduction - Motivation

**Maia**
- **Training partners**: Engines at human skill levels
- **Chess education**: Learn from human-like mistakes
- **Research**: Model human decision-making in games

**Rylee** - all the above AND
- Playing chess on edge devices
- Maia can't play openings

---

# Methodology - Proposed Solution

```txt
+------------------+     +------------------+     +------------------+
|  Human Games     | --> |  Neural Network  | --> |  Move Prediction |
|  (Lichess)       |     |  (CNN + FC)      |     |  (2104 classes)  |
+------------------+     +------------------+     +------------------+
```

Train on millions of human chess positions to predict the next move.

---

# Methodology - Theories

<!-- TODO: Add content about:
     - Theoretical foundations (supervised learning, CNNs for spatial pattern recognition)
     - Why CNNs work for chess (spatial relationships on board)
     - Why neural networks can model human decision-making vs optimal play
     - Trade-offs between model size and performance (why small model can work)
     - Connection to Maia's approach and how you differ
-->

---

# Methodology - ML Models

<!-- TODO: Consider adding visual diagrams:
     - Architecture diagram showing network flow
     - Data pipeline flowchart (more detailed than proposed solution slide)
     - Sample chess board encoding visualization
     - These visuals improve "figure visibility" in grading criteria
-->

<div class="columns">
<div>

### Data Pipeline
- **Download**: .zst files from Lichess
- **Extract**: PGNs
- **Split**: PGNs into individual games
- **FEN snapshots**: Convert to board states
- **Extract metadata**: ELO, result
- **Encode**: One-hot encode snapshots

</div>
<div>

### Neural Network Architecture
- **Input**: Board(12x8x8) + Metadata(4)
- **ConvBlock**: 6x64 @ 8x8, ReLU
- **Flatten**: 4096 + Meta(4)
- **Fully Connected**: 4100 -> 512 -> 32
- **Heads**:
  - MoveHead: 32 -> 2104
  - AuxHead: 32 -> 2104

### Training Settings
- **Chosen Move Loss**: CrossEntropyLoss
- **Valid Moves Loss**: BCE with LogitsLoss
- **Optimizer**: Adam
- **Hyperparameters**: learning_rate, decay, beta1, beta2
- **Search Method**: Random Search

</div>
</div>

---

# Experiments - Dataset

<div class="columns">
<div>

* **Source:** Lichess Open Database
* **Games:** Human rated games
* **Positions:** Each game produces multiple board snapshots (one per move)
* **Scope:** 860,000 million games -> 25 million snapshots
* **Action Space:** 2,104 legal move classes
* **Legal Moves**: Indexes of Legal Moves from a given board state

</div>
<div>

| Split      | Percentage | Snapshots      |
|------------|------------|----------------|
| Training   | 80%        | **20,000,000** |
| Validation | 10%        | **2,500,000**  |
| Test       | 10%        | **2,500,000**  |

</div>
</div>

---

# Experiments - Baselines

| Method            | Description                                         |
|-------------------|-----------------------------------------------------|
| **Random**        | Uniform random choice among all legal moves         |
| **Random Forest** | Classic ML baseline using handcrafted features      |
| **Stockfish 15**  | Strong traditional engine (anchor baseline)         |
| **Leela 4200**    | High-strength neural chess engine (anchor baseline) |
| **Maia1 1500**    | Human-aligned prediction model (anchor baseline)    |

---

# Experiments - Evaluation Metrics

<!-- TODO: Add more evaluation metrics:
     - Consider adding: perplexity, move ranking metrics, or accuracy by game phase
     - Explain why these metrics are appropriate for human move prediction
-->

- **Top-1 Accuracy**: Predicted move is the actual move
- **Top-5 Accuracy**: Predicted move is one of the top predicted moves

# Experiments - Comparisons

<!-- TODO: Consider enhancing comparisons:
     - Add visualization (bar chart or graph comparing methods)
     - Include error bars or confidence intervals if available
     - Add more metrics beyond Top-1 (Top-5, Top-10?)
     - Discuss statistical significance of differences
     - Add comparison of model sizes/computational requirements
-->

| Method            | Top-1 Accuracy |
|-------------------|----------------|
| **Random**        | **6%**         |
| **Random Forest** | **13%**        |
| **Stockfish 15**  | **40%**        |
| **Leela 4200**    | **44%**        |
| **Maia1 1500**    | **51%**        |
| **Rylee (Ours)**  | **20%**        |

---

# Conclusions - Discussions

<!-- TODO:
     - Add discussion of what these numbers mean
     - Discuss why Rylee achieves 20% vs Maia's 51%
     - Discuss model size trade-offs (Raspberry Pi vs full GPU)
     - Add insights about what worked and what didn't
     - Consider adding: training time, model size comparison, inference speed
-->

| Metric          | Value  |
|-----------------|--------|
| Training Loss   | 0.0152 |
| Training Top-1 Accuracy  | 21.8578 |
| Training Top-5 Accuracy  | 47.1662 |
| Validation Loss | 0.0164 |
| Validation Top-1 Accuracy  | 19.4929 |
| Validation Top-5 Accuracy  | 43.1566 |

---

# Conclusions - Future Work

- More training data
- More computational resources
- Larger models
- Data augmentation (board flipping)
- Elo predictor
- Human vs Bot discriminator
- GAN implementation
- Blunder detection

---

# References

<!-- TODO: ADD REFERENCES SLIDE HERE:
     Required for "references used" in grading criteria. Include:
     - Maia Chess: https://maiachess.com/ and the Maia paper
     - Lichess Open Database: https://database.lichess.org/
     - Stockfish: https://stockfishchess.org/
     - Leela Chess Zero: https://lczero.org/
     - Any papers, frameworks, or tools you referenced
     Format as a proper references slide with citations
-->

---

<!-- _class: lead -->

# Questions?

**Rylee**: Human-Like Chess Engine

Ethan Gee & Nate Stott
