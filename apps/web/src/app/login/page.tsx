import LoginForm from "./login-form";

export default async function LoginPage({
  searchParams,
}: {
  searchParams: Promise<{ next?: string }>;
}) {
  const { next } = await searchParams;

  return (
    <div className="flex flex-1 items-center justify-center bg-zinc-50 px-6 dark:bg-black">
      <div className="w-full max-w-sm rounded-lg border border-zinc-200 p-8 dark:border-zinc-800">
        <h1 className="text-xl font-semibold">Sign in to Vakalat</h1>
        <p className="mt-1 text-sm text-zinc-500">
          Enter your credentials to access your matters.
        </p>
        <LoginForm next={next ?? "/"} />
      </div>
    </div>
  );
}
