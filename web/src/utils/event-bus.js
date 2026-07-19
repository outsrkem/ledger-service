// event-bus.js
import mitt from "mitt";

const EventBusPlugin = {
    install(app) {
        app.config.globalProperties.$globalBus = mitt();
    },
};
export default EventBusPlugin;
