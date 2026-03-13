import { PAGES } from "@constants/page-routes";
import { AppLayout, ButtonGroup, FormGroup } from "@versaur/react/blocks";
import { EmailInput, PasswordInput, TextInput } from "@versaur/react/forms";
import { Button } from "@versaur/react/primitive";

export const SignUpPage = () => {
  return (
    <AppLayout>
      <AppLayout.Body centered>
        <AppLayout.Main>
          <h1>Create Account</h1>
          <FormGroup>
            <FormGroup.Field>
              <TextInput placeholder="Full Name" label="Full Name" required />
            </FormGroup.Field>
            <FormGroup.Field>
              <EmailInput placeholder="Email" label="Email" required />
            </FormGroup.Field>
            <FormGroup.Field>
              <PasswordInput placeholder="Password" label="Password" required />
            </FormGroup.Field>
            <FormGroup.Field>
              <PasswordInput placeholder="Confirm Password" label="Confirm Password" required />
            </FormGroup.Field>
            <FormGroup.Field>
              <ButtonGroup>
                <Button type="submit">Sign Up</Button>
              </ButtonGroup>
            </FormGroup.Field>
            <FormGroup.Field>
              <span>
                Already have an account? <a href={PAGES.SIGN_IN}>Sign In</a>
              </span>
            </FormGroup.Field>
          </FormGroup>
        </AppLayout.Main>
      </AppLayout.Body>
    </AppLayout>
  );
};
