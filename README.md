# Backoff

This package implements concurrency safe and reusable backoff algorithms. All the algorithms will have the optional
jitter applied after the delay has been calculated.

Concurrency safety relies on the struct fields not being changed when they might be in use. Once set, the code will not
make any changes to the fields.

## Algorithms

| Algorithm   | Formula           | Description                                                                                                                                  |
|-------------|-------------------|----------------------------------------------------------------------------------------------------------------------------------------------|
| Constant    | `x(n) = b`        | Causes a constant delay between each adverse event.                                                                                          |
| Linear      | `x(n) = b × a`    | Causes a linearly increasing backoff delay that scales with the number of adverse events.                                                    |
| Exponential | `x(n) = b × fⁿ⁻¹` | Causes an exponentially increasing backoff delay that scales with the number of adverse events and the factor (`f`) of at least 2 or higher. |

The optional jitter (`j`) is added to the delay returned by the algorithms. A constant delay with jitter can be
expressed as `x(n) = n ± j`.

## What is it for?

Backoff algorithms are used to add a delay between attempts of an operation that can fail but can be repeated.
An example is connecting to a database where one failed attempt should not cause the application to just exit. Repeating
attempts should not be done at the highest possible speed that an application can do it since it would just cause
unnecessary CPU usage and network traffic that can raise alarms in firewalls.
