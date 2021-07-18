import React from 'react';
import paramsSchema from "./schema-params.json"
import connectionsSchema from "./schema-connections.json"
import artifactsSchema from "./schema-artifacts.json"
import { withTheme } from '@rjsf/core'
import { Theme as MaterialUITheme } from '@rjsf/material-ui'

const type = "test-bundle"
const Form = withTheme(MaterialUITheme)

const log = (type) => console.log.bind(console, type);
const ParamsTemplate = (args) => <Form schema={paramsSchema} onChange={log("changed")} onSubmit={log("submitted")} onError={log("errors")} />;
const ConnectionsTemplate = (args) => <Form schema={connectionsSchema} onChange={log("changed")} onSubmit={log("submitted")} onError={log("errors")} />;
const ArtifactsTemplate = (args) => <Form schema={artifactsSchema} onChange={log("changed")} onSubmit={log("submitted")} onError={log("errors")} />;

export default {
  title: `Bundles/${type}`,
  component: ParamsTemplate,
};

export const Params = ParamsTemplate.bind({});
Params.args = {};

export const Connections = ConnectionsTemplate.bind({});
Connections.args = {};

export const Artifacts = ArtifactsTemplate.bind({});
Artifacts.args = {};
