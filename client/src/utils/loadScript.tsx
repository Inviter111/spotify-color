interface IScriptAttributes {
    async?: boolean;
    defer?: boolean;
    id?: string;
    source: string;
  }

export function loadScript(attributes: IScriptAttributes): Promise<any> {
  if (!attributes || !attributes.source) {
    throw new Error("Invalid attributes");
  }

  return new Promise((resolve, reject) => {
    const { async, defer, id, source }: IScriptAttributes = {
      async: false,
      defer: false,
      source: "",
      ...attributes
    };

    const scriptTag = document.getElementById("spotify-player");

    if (!scriptTag) {
      const script = document.createElement("script");

      script.id = id || "";
      script.type = "text/javascript";
      script.async = async;
      script.defer = defer;
      script.src = source;
      script.onload = () => resolve(undefined);
      script.onerror = (error: any) => reject(`createScript: ${error.message}`);

      document.head.appendChild(script);
    } else {
      resolve();
    }
  });
}
