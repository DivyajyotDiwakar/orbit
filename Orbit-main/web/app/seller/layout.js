import { requireUser } from "@/lib/auth";
export default async function SellerLayout({ children }) {
  await requireUser("seller");
  return children;
}
