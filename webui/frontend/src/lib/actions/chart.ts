export function chart(node: HTMLElement, options: any) {
    let chart: any;

    async function init() {
        // Dynamic import to avoid SSR issues
        const ApexCharts = (await import('apexcharts')).default;
        chart = new ApexCharts(node, options);
        chart.render();
    }

    init();

    return {
        update(newOptions: any) {
            if (chart) {
                chart.updateOptions(newOptions);
            }
        },
        destroy() {
            if (chart) {
                chart.destroy();
            }
        }
    };
}
