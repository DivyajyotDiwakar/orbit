"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";

export function SignupForm(props) {
  const router = useRouter();

  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [age, setAge] = useState("");
  const [role, setRole] = useState("");

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();

    const trimmedFirstName = firstName.trim();
    const trimmedLastName = lastName.trim();
    const trimmedEmail = email.trim().toLowerCase();

    if (!trimmedFirstName || !trimmedLastName) {
      alert("Please enter your full name.");
      return;
    }

    if (!trimmedEmail) {
      alert("Please enter your email.");
      return;
    }

    if (!role) {
      alert("Please choose how you want to use Orbit.");
      return;
    }

    if (password !== confirmPassword) {
      alert("Passwords do not match.");
      return;
    }

    try {
      const res = await fetch("/api/auth/signup", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({
          firstName: firstName,
          lastName: lastName,
          emailId: trimmedEmail,
          password,
          age: Number(age),
          role,
        }),
      });

      if (res.ok) {
        router.push("/");
        return;
      }

      const err = await res.json();
      console.error(err);
      alert(err.message || "Signup failed.");
    } catch (err) {
      console.error(err);
      alert("Something went wrong.");
    }
  };

  return (
    <Card {...props}>
      <CardHeader>
        <CardTitle>Create an account</CardTitle>
        <CardDescription>
          Enter your information below to create your account.
        </CardDescription>
      </CardHeader>

      <CardContent>
        <form onSubmit={handleSubmit}>
          <FieldGroup>
            <Field>
              <FieldLabel htmlFor="firstName">First Name</FieldLabel>

              <Input
                id="firstName"
                placeholder="John"
                value={firstName}
                onChange={(e) => setFirstName(e.target.value)}
                required
              />
            </Field>

            <Field>
              <FieldLabel htmlFor="lastName">Last Name</FieldLabel>

              <Input
                id="lastName"
                placeholder="Doe"
                value={lastName}
                onChange={(e) => setLastName(e.target.value)}
                required
              />
            </Field>

            <Field>
              <FieldLabel htmlFor="age">Age</FieldLabel>

              <Input
                id="age"
                type="number"
                min={18}
                placeholder="18"
                value={age}
                onChange={(e) => setAge(e.target.value)}
                required
              />

              <FieldDescription>
                You must be at least 18 years old.
              </FieldDescription>
            </Field>

            <Field>
              <FieldLabel>Email</FieldLabel>

              <Input
                type="email"
                placeholder="john@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />

              <FieldDescription>
                We&apos;ll use this email to contact you.
              </FieldDescription>
            </Field>

            <Field>
              <FieldLabel>Password</FieldLabel>

              <Input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />

              <FieldDescription>
                Must be at least 8 characters long.
              </FieldDescription>
            </Field>

            <Field>
              <FieldLabel>Confirm Password</FieldLabel>

              <Input
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                required
              />
            </Field>

            <Field>
              <FieldLabel>How will you use Orbit?</FieldLabel>

              <div className="grid grid-cols-2 gap-3">
                <Button
                  type="button"
                  variant={role === "buyer" ? "default" : "outline"}
                  onClick={() => setRole("buyer")}
                >
                  I want to buy
                </Button>

                <Button
                  type="button"
                  variant={role === "seller" ? "default" : "outline"}
                  onClick={() => setRole("seller")}
                >
                  I want to sell
                </Button>
              </div>

              <FieldDescription>
                Your selected account role is used by the existing Orbit
                backend.
              </FieldDescription>
            </Field>

            <Button type="submit" className="w-full">
              Create Account
            </Button>

            <FieldDescription className="text-center">
              Already have an account?{" "}
              <Link
                href="/signin"
                className="font-medium underline underline-offset-4"
              >
                Sign in
              </Link>
            </FieldDescription>
          </FieldGroup>
        </form>
      </CardContent>
    </Card>
  );
}
