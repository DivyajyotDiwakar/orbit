import { requireUser } from "@/lib/auth";
export default async function AdminLayout({ children }) {
  await requireUser("admin");
  return children;
}
