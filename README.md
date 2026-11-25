# BitNet Runner

[![Release](https://shields.io)](https://github.com/mibrahimzia/bitnet-runner/releases)
[![License](https://img.shields.io/github/license/mibrahimzia/bitnet-runner)](https://github.com)

A lightweight desktop app for offline BitNet model inference, free and open source.

## Quick Start

1. **Install the App**  
   Download and install BitNet Runner for Windows from the [Releases page](https://github.com/mibrahimzia/bitnet-runner/releases). On first launch, it will set up a `.bitnet-runner` folder under your user directory.

2. **Add a Model**  
   Place any BitNet model in `.gguf` format inside the `models` folder located at:
   ```text
   C:\Users\<username>\.bitnet-runner\models\
   ```
   Restart the app if it is already running so the model appears in the dropdown menu.

3. **Start Chatting Offline**  
   Launch the app, select your model, and begin chatting through the built-in interface.

👉 **Example model to try:** [microsoft/bitnet-b1.58-2B-4T-gguf](https://huggingface.co/microsoft/bitnet-b1.58-2B-4T-gguf)

---

## Usage Instructions

### 1. Install the Desktop App
Download and run the Windows installer (`BitNetRunner_Setup_v1.0.0`) from the [Releases page](https://github.com/mibrahimzia/bitnet-runner/releases). After installation, the app creates a working directory at:

```text
C:\Users\<username>\.bitnet-runner\
```

### 2. Add Models
Inside the `.bitnet-runner` directory, open the `models` folder:

```text

C:\Users\<username>\.bitnet-runner\(\models\ \%\%\)MAGIT_PARSER_PROTECT%%```

Place your BitNet models here in **`.gguf` format**. If the app is already running, close and restart it so the new model is detected and shown in the dropdown menu.

### 3. Download a Model from Hugging Face
To get started quickly, download Microsoft's official 1.58-bit 2B model: [BitNet b1.58 2B 4T GGUF](https://huggingface.co/microsoft/bitnet-b1.58-2B-4T-gguf)

1. Visit the model page on Hugging Face.
2. Click the **Files and versions** tab.
3. Download the `.gguf` file and move it into your `models` folder.

---
```
## About the Project

BitNet Runner wraps `bitnet.cpp` into a plug-and-play desktop app to make testing 1.58-bit BitNet models simple and stress-free—without manually compiling C++ toolchains or managing dependencies.

BitNet is a native 1.58-bit LLM architecture designed for fast, low-memory, energy-efficient CPU inference:

<img width="2313" height="1235" alt="bitnet_graph" src="https://github.com/user-attachments/assets/78f94017-3bae-4bf3-8a38-579bdc5d94df" />
