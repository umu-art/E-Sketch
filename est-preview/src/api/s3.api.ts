import { GetObjectCommand, S3Client } from '@aws-sdk/client-s3';
import { Upload } from '@aws-sdk/lib-storage';

const s3Client = new S3Client({
  region: process.env.S3_REGION || 'ru-1',
  endpoint: process.env.S3_ENDPOINT,
  credentials: {
    accessKeyId: process.env.S3_ACCESS_KEY_ID,
    secretAccessKey: process.env.S3_SECRET_ACCESS_KEY,
  },
  forcePathStyle: true,
});

const BUCKET_NAME = process.env.S3_BUCKET_NAME;

export async function save(boardId: string, preview: Buffer): Promise<void> {
  const upload = new Upload({
    client: s3Client,
    params: {
      Bucket: BUCKET_NAME,
      Key: `previews/${boardId}.jpeg`,
      Body: preview,
      ContentType: 'image/jpeg',
    },
  });

  try {
    await upload.done();
    console.log(`Preview for board ${boardId} saved successfully.`);
  } catch (error) {
    console.error(`Error saving preview for board ${boardId}:`, error);
    throw error;
  }
}

export async function load(boardId: string): Promise<Buffer | undefined> {
  const getObjectParams = {
    Bucket: BUCKET_NAME,
    Key: `previews/${boardId}.jpeg`,
  };

  try {
    const { Body } = await s3Client.send(new GetObjectCommand(getObjectParams));
    if (Body) {
      const bodyContents = await streamToBuffer(Body as NodeJS.ReadableStream);
      console.log(`Preview for board ${boardId} loaded successfully.`);
      return bodyContents;
    }
  } catch (error) {
    if ((error as any).name === 'NoSuchKey') {
      console.log(`No preview found for board ${boardId}.`);
      return undefined;
    }
    console.error(`Error loading preview for board ${boardId}:`, error);
    throw error;
  }
}

async function streamToBuffer(stream: NodeJS.ReadableStream): Promise<Buffer> {
  return new Promise((resolve, reject) => {
    const chunks: any[] = [];
    stream.on('data', (chunk) => chunks.push(chunk));
    stream.on('error', reject);
    stream.on('end', () => resolve(Buffer.concat(chunks)));
  });
}