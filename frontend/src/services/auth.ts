import {
	getAccessToken,
	getRefreshToken,
	isJwtTokenValid,
	removeRefreshToken,
	setRefreshToken,
} from "$lib/auth";
import { API_URL } from "$lib/consts";
import { fetchAuth } from "$lib/fetchAuth";

export const refreshAccessToken = async (): Promise<string | null> => {
	try {
		const refreshToken = getRefreshToken();
		if (!isJwtTokenValid(refreshToken)) {
			removeRefreshToken();
			throw "Refresh token is invalid";
		}

		const res = await fetch(`${API_URL}/auth/refresh-token`, {
			method: "POST",
			body: JSON.stringify({
				refreshToken: refreshToken,
			}),
		});

		if (!res.ok) {
			removeRefreshToken();
			throw "Refresh token is invalid";
		}

		const data = (await res.json()) as {
			accessToken: string;
			refreshToken: string;
		};

		setRefreshToken(data.refreshToken);
		return data.accessToken;
	} catch (e) {
		console.log("refreshAccessToken: ", e);
	}
	return null;
};

export const fetchAccessToken = async (): Promise<string | null> => {
	const tokenInMemory = getAccessToken();
	if (isJwtTokenValid(tokenInMemory)) {
		return tokenInMemory;
	}

	return await refreshAccessToken();
};

interface ISignInResult {
	userId: number;
	accessToken: string;
}

export const signInWithEmail = async (
	email: string,
	password: string,
): Promise<ISignInResult | null> => {
	try {
		const res = await fetch(`${API_URL}/auth/sign-in`, {
			method: "POST",
			body: JSON.stringify({
				email,
				password,
			}),
		});

		const data = (await res.json()) as ISignInResult & { refreshToken: string };

		if (!res.ok) {
			throw data;
		}

		setRefreshToken(data.refreshToken);
		return data;
	} catch (e) {
		console.log("signInWithEmail: ", e);
	}
	return null;
};

export const signUp = async (
	email: string,
	password: string,
): Promise<ISignInResult | null> => {
	try {
		const res = await fetch(`${API_URL}/auth/sign-up`, {
			method: "POST",
			body: JSON.stringify({
				email,
				password,
			}),
		});

		const data = (await res.json()) as ISignInResult & { refreshToken: string };

		if (!res.ok) {
			throw data;
		}

		setRefreshToken(data.refreshToken);
		return data;
	} catch (e) {
		console.log("signUpWithEmail: ", e);
	}
	return null;
};

export const signInWithToken = async (
	token: string,
): Promise<ISignInResult | null> => {
	try {
		const res = await fetch(`${API_URL}/auth/sign-in-with-token`, {
			method: "GET",
			headers: {
				Authorization: `Bearer ${token}`,
			},
		});

		const data = (await res.json()) as ISignInResult;

		if (!res.ok) {
			throw data;
		}
		return { userId: data.userId, accessToken: token };
	} catch (e) {
		console.log("signInWithToken: ", e);
	}
	return null;
};

export const signOut = async (): Promise<boolean> => {
	try {
		fetch;
		const res = await fetchAuth(`${API_URL}/auth/sign-out`, {
			method: "GET",
		});

		const data = (await res.json()) as ISignInResult;

		if (!res.ok) {
			throw data;
		}
		return true;
	} catch (e) {
		console.log("signOut: ", e);
	}
	return false;
};
