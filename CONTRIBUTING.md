# Contributing

The following information provides a set of guidelines for contributing to the Zenrock repo. Use your best judgment, and, if you see room for improvement, please propose changes to this document.

The contributing guide for Zenrock explains the branching structure, how to use the SDK fork, how to make / test updates to SDK branches and how to create release notes.

Contributions come in the form of writing documentation, raising issues / PRs, and any other actions that help develop the Zenrock documentation.

## First steps

The first step is to find an issue you want to fix. 

If you have a feature request, please [make an issue](https://github.com/Zenrock-Foundation/zrchain/issues/new/choose) for anything of substance, or posting an issue if you want to work on it.

Once you find an existing issue that you want to work on or if you have a new issue to create, continue below.

## Proposing Changes

To contribute a change proposal, use the following workflow:

1. [Fork the repository](https://github.com/Zenrock-Foundation/zrchain/).
2. [Add an upstream](https://docs.github.com/en/github/collaborating-with-pull-requests/working-with-forks/syncing-a-fork) so that you can update your fork.
3. Clone your fork to your computer.
4. Create a branch and name it appropriately.
5. Work on only one major change in one pull request.
6. Make sure all tests are passing locally.
7. Next, rinse and repeat the following:
    1. Commit your changes. Write a simple, straightforward commit message. To learn more, see [How to Write a Git Commit Message](https://chris.beams.io/posts/git-commit/).
    2. Push your changes to your remote fork. To add your remote, you can copy/paste the following:

    ```sh

    #Remove origin

    git remote remove origin

    #set a new remote

    git remote add my_awesome_new_remote_repo [insert-link-found-in-source-subtab-of-your-repo]

    #Verify new remote

    git remote -v

    > my_awesome_new_remote_repo  [link-found-in-source-subtab-of-your-repo] (fetch)
    > my_awesome_new_remote_repo  [link-found-in-source-subtab-of-your-repo] (push)

    #Push changes to your remote repo

    git push <your_remote_name>

    #e.g. git push my_awesome_new_remote_repo
    ```

    3. Create a PR on the Zenrock repository. 
    4. Wait for your changes to be reviewed. If you are a maintainer, you can assign your PR to one or more reviewers. If you aren't a maintainer, one of the maintainers will assign a reviewer.
    5. After you receive feedback from a reviewer, make the requested changes, commit them to your branch, and push them to your remote fork again.
    6. Once approval is given, feel free to squash & merge!

## Git conventions

- Pull Requests are merged using "Squash and merge".
- Pull Requests that are not in draft should have a title that follows
  [conventional commits](https://www.conventionalcommits.org/).
- Keep your Pull Requests as atomic as possible. They should leave the system
  in a working state.

### Commit Messages & Semantic Versioning

We use [`semantic-release`](https://semantic-release.gitbook.io/) to publish releases automatically based on commit history. Every commit that lands on `main` must follow the `type(optional-scope)!: description` Conventional Commits format so the release pipeline can determine the correct semantic version bump.

| Commit type / note        | Release impact | When to use |
|---------------------------|----------------|-------------|
| `feat`                    | Minor          | New functionality that is backwards compatible |
| `fix`, `perf`             | Patch          | Bug fixes or performance improvements |
| `BREAKING CHANGE` note or `type!` | Major          | Backwards-incompatible changes (document the breaking change in the commit body) |
| `docs`, `chore`, `refactor`, `test`, etc. | No version bump | Internal-only changes that should still follow the Conventional Commits format |

Because PRs are squashed, make sure the final squash commit message (and PR title) also follows this format.

## Scaffolder

We recommend using the `make proto` scaffolder to generate messages, queries, types and other features. The script builds the proto files as well as the messages, types, and other elements automatically. Before pushing, we highly recommend to run the [init script](./README.md#commands) to make sure the chain builds correctly. 
The script is derived from [ignite](https://docs.ignite.com). For a more in-depth explanation, follow the Ignite CLI's [scaffolder documentation](https://docs.ignite.com/guide/blog/scaffolding).

## Common Security Considerations

There are several security patterns that come up frequently enough to be synthesized into general rules of thumb. While the high level risks are appchain agnostic, the details are mostly tailored to contributing to Zenrock. This is, of course, not even close to a complete list – just a few considerations to keep in mind.

## Panics

It is common for unexpected behavior in Go code to be guardrailed with panics. There are of course times when panics are appropriate to use instead of errors, but it is important to keep in mind that panics in module-executed code will cause the chain to halt.

While these halts are not terribly difficult to recover, they still pose a valid attack vector, especially if the panics can be triggered repeatably.

Thus, we should be cognisant of when we use panics and ensure that we do not panic on behavior that could be very well handled by an error.
