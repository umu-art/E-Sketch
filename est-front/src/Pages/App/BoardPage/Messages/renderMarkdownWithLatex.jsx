import React from 'react';
import ReactMarkdown from 'react-markdown';
import rehypeKatex from 'rehype-katex';
import remarkMath from 'remark-math';
import 'katex/dist/katex.min.css';

const preprocess = (text) => {
  text = text.replace(/[^\S\n]*\\\[/g, '\\[');
  text = text.replace(/[^\S\n]*\\\]/g, '\\]');

  text = text.replace(/\$/g, '\\$');

  text = text.replace(/\\\((.*?)\\\)/g, '$$$1$$');

  text = text.replace(/\\\[([\s\S]*?)\\\]/g, '$$$\n$1\n$$$');

  text = text.replace(/\\frac\{([^]+)\/([^}]+)\}/g, '\\frac{$1}{$2}');

  return text;
}

const renderMarkdownWithLatex = (markdownText) => {
  return (
    <ReactMarkdown
      remarkPlugins={[remarkMath]} // Подключаем remarkMath для распознавания математики
      rehypePlugins={[rehypeKatex]} // Подключаем rehypeKatex для рендеринга математики
    >
      {preprocess(markdownText)}
    </ReactMarkdown>
  );
};

export default renderMarkdownWithLatex;