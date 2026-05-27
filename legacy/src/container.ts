import {DIContainer} from "@wessberg/di";
import Validator from "./validator/index.js";

const container = new DIContainer();

container.registerSingleton<Validator>()

export default container
