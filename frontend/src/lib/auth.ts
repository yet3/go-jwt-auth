import { jwtDecode } from "jwt-decode";

let accessToken: string | null = null;

export const getAccessToken = () => {
  return accessToken;
}

export const setAccessToken = (token: string | null) => {
	accessToken = token;
};

export const getRefreshToken = (): string | null => {
	return localStorage.getItem("refresh_token");
};

export const setRefreshToken = (token: string) => {
	localStorage.setItem("refresh_token", token);
};

export const removeRefreshToken = () => {
	localStorage.removeItem("refresh_token");
};

interface JwtPayload {
	userId: number;
	exp: number;
	iat: number;
}

// Just check whether it is expired or an valid jwt token
export const isJwtTokenValid = (
	token: string | null | undefined,
): boolean => {
	if (!token) return false;

	const { exp } = jwtDecode<JwtPayload>(token);
	if (new Date().getTime() / 1000 >= exp) {
		return false;
	}

	return true;
};
