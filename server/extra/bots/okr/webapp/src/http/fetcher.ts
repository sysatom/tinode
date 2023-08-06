export default async function fetcher<JSON = any>(
  action: string,
  content?: Object | String
): Promise<JSON> {
  // @ts-ignore
  const bearerToken = Global.flag;

  if (!bearerToken) {
    throw new Error("No token found");
  }
  const data = {
    "action": action,
    "version": 1,
    "content": content
  }
  const res = await fetch("", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${bearerToken}`
    },
    body: JSON.stringify(data)
  });

  if (!res.ok) {
    throw new Error("Not authenticated");
  }

  return res.json();
}
