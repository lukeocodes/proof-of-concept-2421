import { config } from "dotenv";
import { createClient } from "@deepgram/sdk";
import fs from "fs";
import express from "express";
import yargs from "yargs";
import { hideBin } from "yargs/helpers";

// Add TypeScript types
import type { Request, Response } from "express";

config();

const DEEPGRAM_API_KEY = process.env.DEEPGRAM_API_KEY;
const DEFAULT_AUDIO_URL = "https://dpgr.am/spacewalk.wav";

if (!DEEPGRAM_API_KEY) {
  console.error("Missing DEEPGRAM_API_KEY in .env");
  process.exit(1);
}

const deepgram = createClient(DEEPGRAM_API_KEY);

interface TranscriptionResponse {
  transcript: string;
}

async function transcribeFromUrl(url: string): Promise<TranscriptionResponse> {
  try {
    const { result, error } = await deepgram.listen.prerecorded.transcribeUrl(
      { url },
      { model: "nova", punctuate: true }
    );

    if (error) {
      throw error;
    }

    return { 
      transcript: result.results?.channels[0]?.alternatives[0]?.transcript || "" 
    };
  } catch (error) {
    console.error("Error transcribing audio from URL:", error);
    throw error;
  }
}

async function transcribeFromFile(filepath: string): Promise<TranscriptionResponse> {
  try {
    const audioFile = fs.readFileSync(filepath);
    const { result, error } = await deepgram.listen.prerecorded.transcribeFile(
      audioFile,
      { model: "nova", punctuate: true }
    );

    if (error) {
      throw error;
    }

    return { 
      transcript: result.results?.channels[0]?.alternatives[0]?.transcript || "" 
    };
  } catch (error) {
    console.error("Error transcribing audio from file:", error);
    throw error;
  }
}

const argv = yargs(hideBin(process.argv))
  .option("url", {
    type: "string",
    description: "URL of the audio file to transcribe",
  })
  .option("path", {
    type: "string",
    description: "Local file path of the audio file to transcribe",
  })
  .option("serve", {
    type: "boolean",
    description: "Run as an HTTP server",
  })
  .parseSync();

if (argv.serve) {
  const app = express();
  const port = process.env.PORT || 3000;

  app.use(express.json());

  app.post("/", async (req: Request, res: Response) => {
    try {
      const { url } = req.body;
      
      if (!url) {
        return res.status(400).json({ error: "URL is required in request body" });
      }

      const result = await transcribeFromUrl(url);
      res.json(result);
    } catch (error) {
      res.status(500).json({ error: "Failed to transcribe audio" });
    }
  });

  app.listen(port, () => {
    console.log(`Server listening at http://localhost:${port}`);
  });
} else {
  async function handleCLI() {
    try {
      let result: TranscriptionResponse;
      
      if (argv.path) {
        result = await transcribeFromFile(argv.path);
      } else {
        result = await transcribeFromUrl(argv.url || DEFAULT_AUDIO_URL);
      }
      
      // Pretty print and word wrap for CLI output
      console.log(JSON.stringify(result, null, 2));
    } catch (error) {
      console.error("Transcription failed");
      process.exit(1);
    }
  }

  handleCLI();
}