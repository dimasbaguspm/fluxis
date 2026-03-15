import { DEEP_LINKS } from "@/constants";
import { cx } from "@/lib/cx";
import { useSessionHandler } from "@/providers/session";
import { useLogin } from "@hooks/use-api";
import { vAlignCenter, vFlex, vFlexCol, vGap2, vGap8, vTextCenter } from "@versaur/core/utilities";
import { AppLayout, ButtonGroup, FormGroup } from "@versaur/react/blocks";
import { EmailInput, PasswordInput } from "@versaur/react/forms";
import { Banner, Button, Heading, Text } from "@versaur/react/primitive";
import { useForm } from "react-hook-form";
import { Link, useNavigate } from "react-router";

interface SignInFormInputs {
  email: string;
  password: string;
}

export const SignInPage = () => {
  const navigate = useNavigate();
  const { setTokens } = useSessionHandler();
  const [signIn, err, { isPending }] = useLogin({
    onSuccess: (data) => {
      setTokens(data.accessToken || null, data.refreshToken || null);
      navigate(DEEP_LINKS.DASHBOARD);
    },
  });

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<SignInFormInputs>({
    defaultValues: { email: "", password: "" },
  });

  const onSubmit = async (data: SignInFormInputs) => {
    await signIn(data);
  };

  return (
    <AppLayout>
      <AppLayout.Body centered>
        <AppLayout.Main className={cx(vFlex, vFlexCol, vAlignCenter)}>
          <div
            className={cx(vFlex, vFlexCol, vGap8)}
            style={{ maxWidth: "28rem", paddingTop: "3rem" }}
          >
            <div className={cx(vFlex, vFlexCol, vGap2, vTextCenter)}>
              <Heading as="h1" size="2xl">
                Sign In
              </Heading>
              <Text>Welcome back to Fluxis</Text>
            </div>

            <FormGroup onSubmit={handleSubmit(onSubmit)}>
              {err && (
                <FormGroup.Field>
                  <Banner variant="warning">{err.message}</Banner>
                </FormGroup.Field>
              )}
              <FormGroup.Field>
                <EmailInput
                  placeholder="Email"
                  label="Email"
                  required
                  error={errors.email?.message}
                  {...register("email", {
                    required: "Email is required",
                    pattern: {
                      value: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
                      message: "Invalid email address",
                    },
                  })}
                />
              </FormGroup.Field>
              <FormGroup.Field>
                <PasswordInput
                  placeholder="Password"
                  label="Password"
                  required
                  error={errors.password?.message}
                  {...register("password", {
                    required: "Password is required",
                  })}
                />
              </FormGroup.Field>
              <FormGroup.Field>
                <ButtonGroup>
                  <Button type="submit" loading={isPending}>
                    {isPending ? "Signing in..." : "Sign In"}
                  </Button>
                </ButtonGroup>
              </FormGroup.Field>
              <FormGroup.Field>
                <Text as="span">
                  Doesn't have an account? <Link to={DEEP_LINKS.SIGN_UP}>Sign Up</Link>
                </Text>
              </FormGroup.Field>
            </FormGroup>
          </div>
        </AppLayout.Main>
      </AppLayout.Body>
    </AppLayout>
  );
};
