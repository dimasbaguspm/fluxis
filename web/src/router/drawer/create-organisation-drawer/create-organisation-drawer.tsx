import { ButtonGroup, Drawer } from "@versaur/react/blocks";
import { Button, Text } from "@versaur/react/primitive";

export const CreateOrganisationDrawer = () => {
  return (
    <>
      <Drawer.Header action={<Drawer.CloseButton />}>
        <Drawer.Title>Create Organisation</Drawer.Title>
      </Drawer.Header>
      <Drawer.Body>
        <Text>Product details content here.</Text>
      </Drawer.Body>
      <Drawer.Footer>
        <ButtonGroup fluid>
          <Button variant="ghost">Cancel</Button>
          <Button>Create</Button>
        </ButtonGroup>
      </Drawer.Footer>
    </>
  );
};
