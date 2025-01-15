import type { ICtxWrapper } from "./ctx.types";

export interface IUser {
  id: number
}

export type IUserCtx  = ICtxWrapper<IUser>;
