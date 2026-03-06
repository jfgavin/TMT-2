## Network
- Network is now formed from multiple LIF neurons
- Inputs are designed for scalability - takes agent function handles, not raw data
- Network uses block design methodology, and allows for rapid changes and tuning
- Output produces deterministic spike, which will trigger our sacrifice behaviour

## Sacrifice
- Agents can now request self sacrifice from the server
- This function has been connected to be triggered by the TMT model output

## Visualisation
- Emulator has been extended in the sidebar with a plot window
- Live TMT model data is streamed for each agent, and can be viewed in real-time as the sim progresses
- This is only the output neuron - the inner workings of the model are a "black box", and hard to expose

## Initial results
- So far, only eliminations has been connected to the MS block
- Successful model testing with agents requesting self-sacrifice following many natural eliminations
- Even with drastic changes in tuning parameters `TauRise`, `TauDecay` and network weights, the observed number of self-sacrifices is relatively stable
- This invariability is a very good sign, as it indicates the model will not need very fine tuning, and can handle many inputs and neurons with predictable behaviour

# Requested Changes
- Make TMT model 4 blocks
    - Mortality Salience
    - Worldview Validation
    - Self-esteem
    - Relationship Validation / Proximity to Others

- Possible switch to Gaussian for resource distribution
- Use of Scree plot for optimal k-many neighbourhood selection
