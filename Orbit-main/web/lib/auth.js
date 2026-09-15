import { redirect } from "next/navigation";
import { ApiError } from "@/lib/api";
import { getCurrentUser } from "@/lib/api-server";

export async function requireUser(role) {
  try {
    const user = await getCurrentUser();
    if (role && user.role !== role) redirect("/");
    return user;
  } catch (error) {
    if (error instanceof ApiError && error.status === 401) redirect("/signin");
    throw error;
  }
}
