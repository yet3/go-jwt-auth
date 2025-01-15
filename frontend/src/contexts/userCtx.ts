import type { IUser, IUserCtx } from "$types/user.types";
import { getContext, hasContext, setContext } from "svelte";

const CTX_KEY = Symbol("userCtx");

export const createUserCtx = (ctx: IUserCtx) => {
	setContext(CTX_KEY, ctx);
};

export const getUserCtx = () => {
	return getContext(CTX_KEY) as IUserCtx;
};

export const useUserCtx = () => {
  return getUserCtx().value as IUser
}

export const hasUserCtx = (): boolean => {
	return hasContext(CTX_KEY) && getUserCtx().value != null;
};
