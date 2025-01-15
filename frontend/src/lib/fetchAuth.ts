import { fetchAccessToken } from "$services/auth";

export const fetchAuth = async (
	input: string | URL,
	init?: RequestInit,
): Promise<Response> => {
	const accessToken = await fetchAccessToken();

	const headers: RequestInit["headers"] = {
		...init?.headers,
		Authorization: `Bearer ${accessToken}`,
	};

	return fetch(input, { ...init, headers });
};
