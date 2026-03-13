import { AppLayout, ButtonGroup, FormGroup } from "@versaur/react/blocks";
import { EmailInput, PasswordInput } from "@versaur/react/forms";
import { Button } from "@versaur/react/primitive";

export const LoginPage = () => {
  return (
    <AppLayout>
      <AppLayout.Body centered>
        <AppLayout.Main>
          <FormGroup>
            <FormGroup.Field>
              <EmailInput placeholder="Email" label="Email" required />
            </FormGroup.Field>
            <FormGroup.Field>
              <PasswordInput placeholder="Password" label="Password" required />
            </FormGroup.Field>
            <FormGroup.Field>
              <ButtonGroup>
                <Button type="submit">Submit</Button>
              </ButtonGroup>
            </FormGroup.Field>
          </FormGroup>
        </AppLayout.Main>
      </AppLayout.Body>
    </AppLayout>
  );
};
