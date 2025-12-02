export const useParseText = () => {
  const parseHtml = (value: string): string => {
    if (value) {
      return value.replace(/\n\r?/g, "<br />");
    }
    return "";
  };

  const parseJson = (value: string): string => {
    return value.replace(/\n\r?/g, "\n");
  };

  return {
    parseHtml,
    parseJson,
  };
};
