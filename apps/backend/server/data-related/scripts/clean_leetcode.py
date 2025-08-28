"""
This file contains a script that cleans the Leetcode question dataset and preps it for LLM entry.
"""

import html
import os
import re
from textwrap import dedent

import pandas as pd

CSV_INPUT_PATH = "../data/leetcode_problems.csv"
CSV_OUTPUT_PATH = "../data/leetcode_cleaned.csv"
TEXT_OUTPUT_DIR = "../data/batches"
BATCH_SIZE = 20

# Columns to extract
QUESTION_ID_COL = "frontendQuestionId"
DESCRIPTION_COL = "description"

""" Helpers """


# Not used at the moment but might need in the future
def strip_html(text):
    """Remove HTML tags from the description field."""
    no_tags = re.sub(r"<[^>]+>", "", str(text)).strip()
    return html.unescape(no_tags)


def format_prompt(df_batch):
    """Format a batch into LLM-friendly prompt text."""
    prompt = (
        dedent(
            """
        Please rewrite the following LeetCode-style problem descriptions.
        Keep the same topic and difficulty.
        Make them sound original to avoid copyright infringement and avoid code hints but ensure
        the fundamental logic behind each question remains intact.
        Output the rewritten question descriptions in JSON with the following format:
        ```
        {
            "id": ID of the question as an integer,
            "rewrittenDescription": "Rewritten description of question goes here"
        },
        ```
        Do not include any other supplementary information in your response.
        Include necessary HTML stylings so the text looks good.
        """
        ).strip()
        + "\n\n"
    )

    for i, row in df_batch.iterrows():
        prompt += f"### Problem {i + 1}\nID: {row[QUESTION_ID_COL]}\n{row['clean_description']}\n\n"

    return prompt.strip()


# Main


def main():
    """Main function"""
    print("Loading CSV...")
    df = pd.read_csv(CSV_INPUT_PATH)

    print("Cleaning HTML from descriptions...")
    df["clean_description"] = df[DESCRIPTION_COL]

    print(f"Saving cleaned CSV to: {CSV_OUTPUT_PATH}")
    df[[QUESTION_ID_COL, "clean_description"]].to_csv(CSV_OUTPUT_PATH, index=False)

    os.makedirs(TEXT_OUTPUT_DIR, exist_ok=True)

    print("Exporting batches for LLM entry...")

    for i in range(0, len(df), BATCH_SIZE):
        batch_df = df.iloc[i : i + BATCH_SIZE]
        prompt_text = format_prompt(batch_df)
        batch_file = os.path.join(TEXT_OUTPUT_DIR, f"batch_{i // BATCH_SIZE + 1}.txt")

        with open(batch_file, "w", encoding="utf-8") as f:
            f.write(prompt_text)

        print(f"→ Wrote {batch_file}")

    print("✅ Done.")


if __name__ == "__main__":
    main()
